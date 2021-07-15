package handlers

import (
	"server/database"
	"server/util"

	"github.com/gofiber/fiber/v2"
)


func GetByIdHandler(c *fiber.Ctx) error {
	// Access dbClient
	dbClient := c.Locals("dbClient").(*database.UsersClient)

	id, retrievalError := util.RetrieveGetByIdRequestData(c)
	if retrievalError != nil {
		return util.HandleParsingError(c, retrievalError)
	}

	isSameUser, err := util.IsRequestFromSameUser(c)
	if (fiber.Error{}) != err {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	isAdmin, err2 := util.IsRequestFromAdmin(c)
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err2.Error(),
		})
	}

	if isSameUser || isAdmin {
		user, err := database.GetById(dbClient, id)
		if (fiber.Error{}) != err {
			return c.Status(err.Code).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(user)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Not allowed to access other user resources",
		})
	}
}