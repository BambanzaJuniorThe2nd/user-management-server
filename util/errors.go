package util

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)


var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrEmptyPassword      = errors.New("password can't be empty")
	ErrInvalidAuthToken   = errors.New("invalid auth-token")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("Unauthorized")
)

func HandleParsingError(c *fiber.Ctx, parsingError error) error{
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": parsingError.Error(),
	})
}