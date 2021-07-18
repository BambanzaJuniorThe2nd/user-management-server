package handlers

import (
	"server/database"
	"server/util"

	"github.com/gofiber/fiber/v2"
)

func ResetPasswordHandler(c *fiber.Ctx) error {
	// Access dbClient
	dbClient := c.Locals("dbClient").(*database.UsersClient)

	isAdmin, err := util.IsRequestFromAdmin(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	id, retrievalError := util.RetrieveResetPasswordRequestData(c)
	if retrievalError != nil {
		return util.HandleParsingError(c, retrievalError)
	}

	if isAdmin {
		err := database.ResetUserPassword(dbClient, id)
		if (fiber.Error{}) != err {
			return c.Status(err.Code).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusOK)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": database.ERROR_MESSAGE_ACCESS_RESTRICTED,
		})
	}
}
