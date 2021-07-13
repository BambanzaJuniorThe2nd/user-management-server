package main

import (
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
	client := database.SetupDatabaseClient()

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
