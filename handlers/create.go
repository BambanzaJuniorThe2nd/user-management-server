package handlers

import (
	"server/database"
	"server/util"

	"github.com/gofiber/fiber/v2"
)

func CreateHandler(c *fiber.Ctx) error {
	// Access dbClient
	dbClient := c.Locals("dbClient").(*database.UsersClient)

	isAdmin, err := util.IsRequestFromAdmin(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userDetails, parsingError := util.RetrieveCreateRequestData(c)
	if parsingError != nil {
		return util.HandleParsingError(c, parsingError)
	}

	if isAdmin {
		user, err := database.CreateByAdmin(dbClient, userDetails)
		if (fiber.Error{}) != err {
			return c.Status(err.Code).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(user)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": database.ERROR_MESSAGE_ACCESS_RESTRICTED,
		})
	}
}
