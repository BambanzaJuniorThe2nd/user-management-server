package handlers

import (
	"server/database"
	"server/util"

	"github.com/gofiber/fiber/v2"
)

func DeleteHandler(c *fiber.Ctx) error {
	// Access dbClient
	dbClient := c.Locals("dbClient").(*database.UsersClient)

	isAdmin, err := util.IsRequestFromAdmin(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	id, retrievalError := util.RetrieveDeleteRequestData(c)
	if retrievalError != nil {
		return util.HandleParsingError(c, retrievalError)
	}

	if isAdmin {
		err := database.Delete(dbClient, id)
		if (fiber.Error{}) != err {
			return c.Status(err.Code).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Not allowed to delete user resources",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
