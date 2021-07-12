package middleware

import (
	"server/database"

	"github.com/gofiber/fiber/v2"
)

func AddDatabaseClientToContext(client *database.UsersClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("dbClient", client)
		return c.Next()
	}
}
