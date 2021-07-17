package database

import (
	"context"
	"errors"
	"fmt"
	"server/config"
	"server/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectDB(ctx context.Context, conf config.MongoConfiguration) *mongo.Database {
	connection := options.Client().ApplyURI(conf.Server)

	client, err := mongo.Connect(ctx, connection)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connected!")
	return client.Database(conf.Database)
}

func SetupDatabaseClient() *UsersClient {
	conf := config.GetConfig()
	ctx := context.TODO()

	db := connectDB(ctx, conf.Mongo)
	collection := db.Collection(conf.Mongo.Collection)

	client := &UsersClient{
		Col: collection,
		Ctx: ctx,
	}

	err := createIndices(client)
	if err != nil {
		panic(err)
	}

	err = createDefaultAdmin(client)
	if err != nil {
		panic(err)
	}

	return client
}

func createIndices(usersClient *UsersClient) error {

	// Create an index model for the field: email
	mod := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	// Create the above index on the users collection
	_, err := usersClient.Col.Indexes().CreateOne(usersClient.Ctx, mod)
	if err != nil {
		return err
	}

	return nil
}

func createDefaultAdmin(usersClient *UsersClient) error {
	users, err := GetAll(usersClient)
	if (fiber.Error{}) != err {
		return errors.New(err.Error())
	}

	if len(users) == 0 {
		fmt.Println("Creating default admin")
		user := models.CreateArgs{
			CreateByAdminArgs: models.CreateByAdminArgs{
				Name:      "John Doe",
				Email:     "admin@gmail.com",
				Title:     "Administrator",
				Birthdate: "1970-01-01",
				IsAdmin:   true,
			},
			Password:          "defaultPassword1!",
		}

		_, err := Create(usersClient, user)
		if (fiber.Error{}) != err {
			return errors.New(err.Error())
		}
	}

	return nil
}
