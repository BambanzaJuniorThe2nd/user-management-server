package handlers

import (
	"server/database"
	"server/util"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	// Read in credentials
	creds, parsingError := util.RetrieveLoginRequestData(c)

	if parsingError != nil {
		return util.HandleParsingError(c, parsingError)
	}

	dbClient := c.Locals("dbClient").(*database.UsersClient)
	res, err := database.Login(dbClient, creds)

	// Check whether fields within err
	// are not set to their zero values
	if (fiber.Error{}) != err {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
