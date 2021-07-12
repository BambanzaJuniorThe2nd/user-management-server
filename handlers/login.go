package handlers

import (
	"server/util"
	"server/database"

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
	if err.Message != "" {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)

	// Read in
	// id := c.Params("id")

	// res, err := db.Get(id)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"message": err.Error(),
	// 	})
	// }

	// return c.Status(fiber.StatusOK).JSON(res)
}
