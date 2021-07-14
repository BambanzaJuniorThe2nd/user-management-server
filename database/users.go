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
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersClient struct {
	Ctx context.Context
	Col *mongo.Collection
}

var DEFAULT_PASSWORD = "defaultPassword1"

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
