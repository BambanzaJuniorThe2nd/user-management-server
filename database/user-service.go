package database

import (
	"context"
	"server/models"
	"server/security"
	"server/util"
	"server/validators"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UsersClient struct {
	Ctx context.Context
	Col *mongo.Collection
}

var DEFAULT_PASSWORD = "defaultPassword1@"

func Login(dbClient *UsersClient, args models.LoginArgs) (models.LoginResult, fiber.Error) {
	result := models.LoginResult{}

	validationError := validators.ValidateLoginArgs(args)
	if validationError != nil {
		return result, fiber.Error{Code: fiber.StatusBadRequest, Message: validationError.Error()}
	}

	// Query user with provided email
	user := models.User{}
	query := bson.D{{Key: "email", Value: args.Email}}

	err := dbClient.Col.FindOne(dbClient.Ctx, query).Decode(&user)
	if (err != nil) || (user.Password != "" && !util.CheckPasswordHash(args.Password, user.Password)) {
		return result, fiber.Error{Code: fiber.StatusUnauthorized, Message: ERROR_MESSAGE_LOGIN_FAILED}
	}

	token, tokenError := security.NewToken(&user)
	if tokenError != nil {
		return result, fiber.Error{Code: fiber.StatusInternalServerError, Message: ERROR_MESSAGE_SOMETHING_WENT_WRONG}
	}

	if !user.IsAdmin {
		return result, fiber.Error{Code: fiber.StatusUnauthorized, Message: ERROR_MESSAGE_ACCESS_RESTRICTED}
	}

	result.Token = token
	result.User = GetSafeUser(user)
	return result, fiber.Error{}
}

func CreateByAdmin(dbClient *UsersClient, args models.CreateByAdminArgs) (models.User, fiber.Error) {
	user := models.User{}

	validationError := validators.ValidateCreateByAdminArgs(args)
	if validationError != nil {
		return user, fiber.Error{Code: fiber.StatusBadRequest, Message: validationError.Error()}
	}

	hashedPassword, err := util.HashPassword(DEFAULT_PASSWORD)
	if err != nil {
		return user, fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	// Parse args.CreateByAdminArgs.Birthdate
	birthdate, err := time.Parse(DATE_FORMAT, args.Birthdate)
	if err != nil {
		return user, fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	// Create a User object
	// mostly from args
	user = models.User{
		Name:      args.Name,
		Email:     args.Email,
		Title:     args.Title,
		Birthdate: birthdate,
		Password:  hashedPassword,
		IsAdmin:   args.IsAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := dbClient.Col.InsertOne(dbClient.Ctx, user)
	if err != nil {
		if err.(mongo.WriteException).WriteErrors[0].Code == 11000 {
			return models.User{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: ERROR_MESSAGE_EMAIL_ALREADY_IN_USE}
		}

		return models.User{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: ERROR_MESSAGE_SOMETHING_WENT_WRONG}
	}

	// get the inserted user
	user = models.User{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}

	if err := dbClient.Col.FindOne(dbClient.Ctx, query).Decode(&user); err != nil {
		return models.User{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	return GetSafeUser(user), fiber.Error{}
}

func UpdateByAdmin(dbClient *UsersClient, id primitive.ObjectID, args models.UpdateByAdminArgs) (models.User, fiber.Error) {
	user := models.User{}

	validationError := validators.ValidateUpdateByAdminArgs(args)
	if validationError != nil {
		return user, fiber.Error{Code: fiber.StatusBadRequest, Message: validationError.Error()}
	}

	// Parse args.CreateByAdminArgs.Birthdate
	birthdate, err := time.Parse(DATE_FORMAT, args.CreateByAdminArgs.Birthdate)
	if err != nil {
		return user, fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	updateDoc := bson.D{
		{Key: "name", Value: args.CreateByAdminArgs.Name},
		{Key: "email", Value: args.CreateByAdminArgs.Email},
		{Key: "title", Value: args.CreateByAdminArgs.Title},
		{Key: "birthdate", Value: birthdate},
		{Key: "isAdmin", Value: args.CreateByAdminArgs.IsAdmin},
		{Key: "updatedAt", Value: time.Now()},
	}

	query := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$set", Value: updateDoc},
	}

	err = dbClient.Col.FindOneAndUpdate(dbClient.Ctx, query, update).Err()
	if err != nil {
		if err.(mongo.WriteException).WriteErrors[0].Code == 11000 {
			return models.User{}, fiber.Error{Code: fiber.StatusNotFound, Message: ERROR_MESSAGE_EMAIL_ALREADY_IN_USE}
		} else if err == mongo.ErrNoDocuments {
			return models.User{}, fiber.Error{Code: fiber.StatusNotFound, Message: ERROR_MESSAGE_USER_NOT_FOUND}
		}

		return models.User{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: ERROR_MESSAGE_SOMETHING_WENT_WRONG}
	}

	// get updated data
	user = models.User{}
	dbClient.Col.FindOne(dbClient.Ctx, query).Decode(&user)

	return GetSafeUser(user), fiber.Error{}
}

func Delete(dbClient *UsersClient, id primitive.ObjectID) fiber.Error {
	// find and delete todo
	query := bson.D{{Key: "_id", Value: id}}

	err := dbClient.Col.FindOneAndDelete(dbClient.Ctx, query).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.Error{Code: fiber.StatusNotFound, Message: ERROR_MESSAGE_USER_NOT_FOUND}
		}

		return fiber.Error{Code: fiber.StatusInternalServerError, Message: ERROR_MESSAGE_SOMETHING_WENT_WRONG}
	}

	return fiber.Error{}
}

func GetAll(dbClient *UsersClient) ([]models.User, fiber.Error) {
	query := bson.D{{}}
	projection := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})

	var users []models.User = make([]models.User, 0)
	cursor, err := dbClient.Col.Find(dbClient.Ctx, query, projection)
	if err != nil {
		return users, fiber.Error{Code: fiber.StatusInternalServerError, Message: ERROR_MESSAGE_SOMETHING_WENT_WRONG}
	}

	// iterate the cursor and decode each item into a User
	if err = cursor.All(dbClient.Ctx, &users); err != nil {
		return users, fiber.Error{Code: fiber.StatusInternalServerError, Message: ERROR_MESSAGE_SOMETHING_WENT_WRONG}
	}

	return users, fiber.Error{}
}

func GetById(dbClient *UsersClient, id primitive.ObjectID) (models.User, fiber.Error) {
	user := models.User{}
	query := bson.D{{Key: "_id", Value: id}}

	err := dbClient.Col.FindOne(dbClient.Ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, fiber.Error{Code: fiber.StatusNotFound, Message: ERROR_MESSAGE_USER_NOT_FOUND}
		}

		return user, fiber.Error{Code: fiber.StatusInternalServerError, Message: ERROR_MESSAGE_SOMETHING_WENT_WRONG}
	}

	return GetSafeUser(user), fiber.Error{}
}

func ResetUserPassword(dbClient *UsersClient, id primitive.ObjectID) fiber.Error {
	hashedPassword, err := util.HashPassword(DEFAULT_PASSWORD)
	if err != nil {
		return fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	updateDoc := bson.D{
		{Key: "password", Value: hashedPassword},
		{Key: "updatedAt", Value: time.Now()},
	}

	query := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$set", Value: updateDoc},
	}

	err = dbClient.Col.FindOneAndUpdate(dbClient.Ctx, query, update).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.Error{Code: fiber.StatusNotFound, Message: ERROR_MESSAGE_USER_NOT_FOUND}
		}

		return fiber.Error{Code: fiber.StatusInternalServerError, Message: ERROR_MESSAGE_SOMETHING_WENT_WRONG}
	}

	return fiber.Error{}
}

func ChangePassword(dbClient *UsersClient, id primitive.ObjectID, args models.ChangePasswordArgs) fiber.Error {
	validationError := validators.ValidateChangePasswordArgs(args)
	if validationError != nil {
		return fiber.Error{Code: fiber.StatusBadRequest, Message: validationError.Error()}
	}

	hashedPassword, err := util.HashPassword(args.Password)
	if err != nil {
		return fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	updateDoc := bson.D{
		{Key: "password", Value: hashedPassword},
		{Key: "updatedAt", Value: time.Now()},
	}

	query := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$set", Value: updateDoc},
	}

	err = dbClient.Col.FindOneAndUpdate(dbClient.Ctx, query, update).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.Error{Code: fiber.StatusNotFound, Message: ERROR_MESSAGE_USER_NOT_FOUND}
		}

		return fiber.Error{Code: fiber.StatusInternalServerError, Message: ERROR_MESSAGE_SOMETHING_WENT_WRONG}
	}

	return fiber.Error{}
}

func CreateDefaultAdmin(dbClient *UsersClient, args models.CreateDefaultAdminArgs) fiber.Error {
	user := models.User{}

	validationError := validators.ValidateCreateDefaultAdminArgs(args)
	if validationError != nil {
		return fiber.Error{Code: fiber.StatusBadRequest, Message: validationError.Error()}
	}

	hashedPassword, err := util.HashPassword(args.Password)
	if err != nil {
		return fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	// Parse args.Birthdate
	birthdate, err := time.Parse(DATE_FORMAT, args.Birthdate)
	if err != nil {
		return fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	// Create a User object
	// mostly from args
	user = models.User{
		Name:      args.Name,
		Email:     args.Email,
		Title:     args.Title,
		Birthdate: birthdate,
		IsAdmin:   args.IsAdmin,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := dbClient.Col.InsertOne(dbClient.Ctx, user)
	if err != nil {
		if err.(mongo.WriteException).WriteErrors[0].Code == 11000 {
			return fiber.Error{Code: fiber.StatusInternalServerError, Message: ERROR_MESSAGE_EMAIL_ALREADY_IN_USE}
		}

		return fiber.Error{Code: fiber.StatusInternalServerError, Message: ERROR_MESSAGE_SOMETHING_WENT_WRONG}
	}

	// get the inserted user
	user = models.User{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}

	if err := dbClient.Col.FindOne(dbClient.Ctx, query).Decode(&user); err != nil {
		return fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	return fiber.Error{}
}
