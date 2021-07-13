package handlers

import (
	"server/database"
	"server/util"

	"github.com/gofiber/fiber/v2"
)

func CreateHandler(c *fiber.Ctx) error {
	// Read in user details
	userDetails, parsingError := util.RetrieveCreateRequestData(c)

	if parsingError != nil {
		return util.HandleParsingError(c, parsingError)
	}

	// Access dbClient
	dbClient := c.Locals("dbClient").(*database.UsersClient)

	isAdmin, err := util.IsRequestFromAdmin(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	
	// If request came from admin
	if isAdmin {
		res, err := database.CreateByAdmin(dbClient, userDetails)
	} else {
		res, err := database.Create(dbClient, userDetails)
	}

	// dbClient := c.Locals("dbClient").(*database.UsersClient)
	// res, err := database.Login(dbClient, userDetails)

	// // Check whether fields within err
	// // are not set to their zero values
	// if (fiber.Error{}) != err {
	// 	return c.Status(err.Code).JSON(fiber.Map{
	// 		"message": err.Message,
	// 	})
	// }

	// return c.Status(fiber.StatusOK).JSON(res)
}