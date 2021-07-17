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
		return result, fiber.Error{Code: fiber.StatusUnauthorized, Message: "Login failed"}
	}

	token, tokenError := security.NewToken(&user)
	if tokenError != nil {
		return result, fiber.Error{Code: fiber.StatusInternalServerError, Message: "Something went wrong"}
	}

	result.Token = token
	result.User = util.GetSafeUser(user)
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
	birthdate, err := time.Parse("2006-01-02", args.Birthdate)
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
			return models.User{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: "email already in use"}
		}

		return models.User{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: "Something went wrong"}
	}

	// get the inserted user
	user = models.User{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}

	if err := dbClient.Col.FindOne(dbClient.Ctx, query).Decode(&user); err != nil {
		return models.User{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	return util.GetSafeUser(user), fiber.Error{}
}

func Create(dbClient *UsersClient, args models.CreateArgs) (models.User, fiber.Error) {
	user := models.User{}

	validationError := validators.ValidateCreateArgs(args)
	if validationError != nil {
		return user, fiber.Error{Code: fiber.StatusBadRequest, Message: validationError.Error()}
	}

	hashedPassword, err := util.HashPassword(args.Password)
	if err != nil {
		return user, fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	// Parse args.CreateByAdminArgs.Birthdate
	birthdate, err := time.Parse("2006-01-02", args.CreateByAdminArgs.Birthdate)
	if err != nil {
		return user, fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	// Create a User object
	// mostly from args
	user = models.User{
		Name:      args.CreateByAdminArgs.Name,
		Email:     args.CreateByAdminArgs.Email,
		Title:     args.CreateByAdminArgs.Title,
		Birthdate: birthdate,
		IsAdmin:   args.CreateByAdminArgs.IsAdmin,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := dbClient.Col.InsertOne(dbClient.Ctx, user)
	if err != nil {
		if err.(mongo.WriteException).WriteErrors[0].Code == 11000 {
			return models.User{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: "email already in use"}
		}

		return models.User{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: "Something went wrong"}
	}

	// get the inserted user
	user = models.User{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}

	if err := dbClient.Col.FindOne(dbClient.Ctx, query).Decode(&user); err != nil {
		return models.User{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	return util.GetSafeUser(user), fiber.Error{}
}

func UpdateByAdmin(dbClient *UsersClient, id primitive.ObjectID, args models.UpdateByAdminArgs) (models.User, fiber.Error) {
	user := models.User{}

	validationError := validators.ValidateUpdateByAdminArgs(args)
	if validationError != nil {
		return user, fiber.Error{Code: fiber.StatusBadRequest, Message: validationError.Error()}
	}

	// Parse args.CreateByAdminArgs.Birthdate
	birthdate, err := time.Parse("2006-01-02", args.CreateByAdminArgs.Birthdate)
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
			return models.User{}, fiber.Error{Code: fiber.StatusNotFound, Message: "email already in use"}
		} else if err == mongo.ErrNoDocuments {
			return models.User{}, fiber.Error{Code: fiber.StatusNotFound, Message: "User not found"}
		}

		return models.User{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: "Something went wrong"}
	}

	// get updated data
	user = models.User{}
	dbClient.Col.FindOne(dbClient.Ctx, query).Decode(&user)

	return util.GetSafeUser(user), fiber.Error{}
}

func Update(dbClient *UsersClient, id primitive.ObjectID, args models.UpdateArgs) (models.User, fiber.Error) {
	user := models.User{}

	validationError := validators.ValidateUpdateArgs(args)
	if validationError != nil {
		return user, fiber.Error{Code: fiber.StatusBadRequest, Message: validationError.Error()}
	}

	hashedPassword, err := util.HashPassword(args.Password)
	if err != nil {
		return user, fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	// Parse args.CreateByAdminArgs.Birthdate
	birthdate, err := time.Parse("2006-01-02", args.Birthdate)
	if err != nil {
		return user, fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	updateDoc := bson.D{
		{Key: "name", Value: args.CreateByAdminArgs.Name},
		{Key: "email", Value: args.CreateByAdminArgs.Email},
		{Key: "title", Value: args.CreateByAdminArgs.Title},
		{Key: "birthdate", Value: birthdate},
		{Key: "password", Value: hashedPassword},
		{Key: "updatedAt", Value: time.Now()},
	}

	query := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$set", Value: updateDoc},
	}

	err = dbClient.Col.FindOneAndUpdate(dbClient.Ctx, query, update).Err()
	if err != nil {
		if err.(mongo.WriteException).WriteErrors[0].Code == 11000 {
			return models.User{}, fiber.Error{Code: fiber.StatusNotFound, Message: "email already in use"}
		} else if err == mongo.ErrNoDocuments {
			return models.User{}, fiber.Error{Code: fiber.StatusNotFound, Message: "User not found"}
		}

		return models.User{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: "Something went wrong"}
	}

	// get updated data
	user = models.User{}
	dbClient.Col.FindOne(dbClient.Ctx, query).Decode(&user)

	return util.GetSafeUser(user), fiber.Error{}
}

func Delete(dbClient *UsersClient, id primitive.ObjectID) fiber.Error {
	// find and delete todo
	query := bson.D{{Key: "_id", Value: id}}

	err := dbClient.Col.FindOneAndDelete(dbClient.Ctx, query).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.Error{Code: fiber.StatusNotFound, Message: "User not found"}
		}

		return fiber.Error{Code: fiber.StatusInternalServerError, Message: "Something went wrong"}
	}

	return fiber.Error{}
}

func GetAll(dbClient *UsersClient) ([]models.User, fiber.Error) {
	query := bson.D{{}}
	projection := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})

	var users []models.User = make([]models.User, 0)
	cursor, err := dbClient.Col.Find(dbClient.Ctx, query, projection)
	if err != nil {
		return users, fiber.Error{Code: fiber.StatusInternalServerError, Message: "Something went wrong"}
	}

	// iterate the cursor and decode each item into a User
	if err = cursor.All(dbClient.Ctx, &users); err != nil {
		return users, fiber.Error{Code: fiber.StatusInternalServerError, Message: "Something went wrong"}
	}

	return users, fiber.Error{}
}

func GetById(dbClient *UsersClient, id primitive.ObjectID) (models.User, fiber.Error) {
	user := models.User{}
	query := bson.D{{Key: "_id", Value: id}}

	err := dbClient.Col.FindOne(dbClient.Ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, fiber.Error{Code: fiber.StatusNotFound, Message: "User not found"}
		}

		return user, fiber.Error{Code: fiber.StatusInternalServerError, Message: "Something went wrong"}
	}

	return util.GetSafeUser(user), fiber.Error{}
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
			return fiber.Error{Code: fiber.StatusNotFound, Message: "User not found"}
		}

		return fiber.Error{Code: fiber.StatusInternalServerError, Message: "Something went wrong"}
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
			return fiber.Error{Code: fiber.StatusNotFound, Message: "User not found"}
		}

		return fiber.Error{Code: fiber.StatusInternalServerError, Message: "Something went wrong"}
	}

	return fiber.Error{}
}
