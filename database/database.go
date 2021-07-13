package database

import (
	"context"
	"fmt"
	"server/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(ctx context.Context, conf config.MongoConfiguration) *mongo.Database {
	connection := options.Client().ApplyURI(conf.Server)

	client, err := mongo.Connect(ctx, connection)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connected!")
	return client.Database(conf.Database)
}
