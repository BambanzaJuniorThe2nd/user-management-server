package handlers

import (
	"server/database"
	"server/util"

	"github.com/gofiber/fiber/v2"
)


func GetByToken(c *fiber.Ctx) error {
	// Access dbClient
	dbClient := c.Locals("dbClient").(*database.UsersClient)

	id, err := util.RetrieveIdFromToken(c)
	if (fiber.Error{}) != err {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := database.GetById(dbClient, id)
	if (fiber.Error{}) != err {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
