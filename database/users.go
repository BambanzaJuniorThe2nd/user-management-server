package database

import (
	"context"
	"server/models"
	"server/security"
	"server/util"
	"server/validators"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersClient struct {
	Ctx context.Context
	Col *mongo.Collection
}

type UserServiceInterface interface {
	// Insert(models.Todo) (models.Todo, error)
	// Update(string, interface{}) (models.TodoUpdate, error)
	// Delete(string) (models.TodoDelete, error)
	// Get(string) (models.Todo, error)
	// Search(interface{}) ([]models.Todo, error)
	Login(models.LoginArgs) (models.User, fiber.Error)
}

func Login(dbClient *UsersClient, args models.LoginArgs) (models.LoginResult, fiber.Error) {
	result := models.LoginResult{}

	validationError := validators.ValidateLoginArgs(args)
	if validationError != nil {
		return result, fiber.Error{Code: fiber.StatusBadRequest, Message: validationError.Error()}
	}

	// Query user with provided email
	user := &models.User{}
	query := bson.D{{Key: "email", Value: args.Email}}

	err := dbClient.Col.FindOne(dbClient.Ctx, query).Decode(user)
	if (err != nil) || (user.Password != "" && !util.CheckPasswordHash(args.Password, user.Password)) {
		return result, fiber.Error{Code: fiber.StatusUnauthorized, Message: "Login failed"}
	}

	token, tokenError := security.NewToken(user)
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
}

// func Create(dbClient *UsersClient, args models.CreateArgs) (models.User, fiber.Error) {
// 	user := models.User{}

// 	validationError := validators.ValidateCreateArgs(args)
// 	if validationError != nil {
// 		return user, fiber.Error{Code: fiber.StatusBadRequest, Message: validationError.Error()}
// 	}
// }
