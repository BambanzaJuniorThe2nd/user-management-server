package main

// import (
// 	"log"
// 	"os"
// 	"server/config"
// 	"server/routes"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/cors"
// 	"github.com/gofiber/fiber/v2/middleware/logger"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	if os.Getenv("APP_ENV") != "production" {
// 		err := godotenv.Load()
// 		if err != nil {
// 			log.Fatal("Error loading .env file")
// 		}
// 	}

// 	app := fiber.New()

// 	app.Use(cors.New())
// 	app.Use(logger.New())

// 	config.ConnectDB()

// 	setupRoutes(app)

// 	port := os.Getenv("PORT")
// 	err := app.Listen(":" + port)

// 	if err != nil {
// 		log.Fatal("Error app failed to start")
// 		panic(err)
// 	}
// }

import (
	"context"
	"log"
	"server/config"
	"server/database"
	"server/middleware"
	"server/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "You are at the root endpoint"})
	})

	api := app.Group("/api")

	routes.UsersRoute(api.Group("/users"))
}

func main() {
	conf := config.GetConfig()
	ctx := context.TODO()

	db := database.ConnectDB(ctx, conf.Mongo)
	collection := db.Collection(conf.Mongo.Collection)

	client := &database.UsersClient{
		Col: collection,
		Ctx: ctx,
	}

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(middleware.AddDatabaseClientToContext(client))

	setupRoutes(app)

	err := app.Listen(":" + config.GetConfig().Port)

	if err != nil {
		log.Fatal("Error app failed to start")
		panic(err)
	}
}
