package handlers

import (
	"server/database"
	"server/models"
	"server/util"

	"github.com/gofiber/fiber/v2"
)

func UpdateHandler(c *fiber.Ctx) error {
	// Access dbClient
	dbClient := c.Locals("dbClient").(*database.UsersClient)

	isAdmin, err := util.IsRequestFromAdmin(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	id, updateData, retrievalError := util.RetrieveUpdateRequestData(c, isAdmin)
	if retrievalError != nil {
		return util.HandleParsingError(c, retrievalError)
	}

	var res models.User
	if isAdmin {
		user, err := database.UpdateByAdmin(dbClient, id, updateData.(models.UpdateByAdminArgs))
		if (fiber.Error{}) != err {
			return c.Status(err.Code).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		res = user
	} else {
		user, err := database.Update(dbClient, id, updateData.(models.UpdateArgs))
		if (fiber.Error{}) != err {
			return c.Status(err.Code).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		res = user
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}
