package handlers

import (
	"fmt"
	"server/database"
	"server/models"
	"server/util"

	"github.com/gofiber/fiber/v2"
)

func GetAllHandler(c *fiber.Ctx) error {
	fmt.Println("Inside GetAllHandler...")
	
	// Access dbClient
	dbClient := c.Locals("dbClient").(*database.UsersClient)

	isAdmin, err := util.IsRequestFromAdmin(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var res []models.User
	if isAdmin {
		users, err := database.GetAll(dbClient)
		if (fiber.Error{}) != err {
			return c.Status(err.Code).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		res = users
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Not allowed to access other user resources",
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
