package handlers

import (
	"server/database"
	"server/models"
	"server/util"

	"github.com/gofiber/fiber/v2"
)

func CreateHandler(c *fiber.Ctx) error {
	// Access dbClient
	dbClient := c.Locals("dbClient").(*database.UsersClient)

	// isAdmin, err := util.IsRequestFromAdmin(c)
	isAdmin := false
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"message": err.Error(),
	// 	})
	// }

	userDetails, parsingError := util.RetrieveCreateRequestData(c, isAdmin)
	if parsingError != nil {
		return util.HandleParsingError(c, parsingError)
	}

	var res models.User
	if isAdmin {
		user, err := database.CreateByAdmin(dbClient, userDetails.(models.CreateByAdminArgs))
		if (fiber.Error{}) != err {
			return c.Status(err.Code).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		res = user
	} else {
		user, err := database.Create(dbClient, userDetails.(models.CreateArgs))
		if (fiber.Error{}) != err {
			return c.Status(err.Code).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		res = user
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}
