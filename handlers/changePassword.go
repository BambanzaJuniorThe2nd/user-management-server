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

	isSameUser, err := util.IsRequestFromSameUser(c)
	if (fiber.Error{}) != err {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if isSameUser {
		err = database.ChangePassword(dbClient, id, args)
		if (fiber.Error{}) != err {
			return c.Status(err.Code).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusOK)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Not authorized to change other users' passwords",
		})
	}
}