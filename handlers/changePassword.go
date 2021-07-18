package handlers

import (
	"server/database"
	"server/util"

	"github.com/gofiber/fiber/v2")


func ChangePasswordHandler(c *fiber.Ctx) error {
	// Access dbClient
	dbClient := c.Locals("dbClient").(*database.UsersClient)

	id, args, err := util.RetrieveChangePasswordRequestData(c)
	if (fiber.Error{}) != err {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	isAdmin, err2 := util.IsRequestFromAdmin(c)
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if isAdmin {
		err = database.ChangePassword(dbClient, id, args)
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