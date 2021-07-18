package util

import (
	"server/models"
	"server/security"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type JError struct {
	Error string `json:"error"`
}

func NewJError(err error) JError {
	jerr := JError{"generic error"}
	if err != nil {
		jerr.Error = err.Error()
	}
	return jerr
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ExtractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")

	onlyToken := strings.Split(bearToken, " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

func ConvertStringIdIntoObjectId(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func RetrieveLoginRequestData(c *fiber.Ctx) (models.LoginArgs, error) {
	// Create an empty creds object
	// with the LoginArgs structure
	creds := models.LoginArgs{}

	err := c.BodyParser(&creds)
	return creds, err
}

func RetrieveCreateRequestData(c *fiber.Ctx) (models.CreateByAdminArgs, error) {
	data := models.CreateByAdminArgs{}
	err := c.BodyParser(&data)
	return data, err
}

func RetrieveUpdateRequestData(c *fiber.Ctx) (primitive.ObjectID, models.UpdateByAdminArgs, error) {
	// Convert id parameter to objectId
	id, err := ConvertStringIdIntoObjectId(c.Params("id"))
	if err != nil {
		return primitive.ObjectID{}, models.UpdateByAdminArgs{}, err
	}

	data := models.UpdateByAdminArgs{}
	err = c.BodyParser(&data)
	return id, data, err
}

func RetrieveDeleteRequestData(c *fiber.Ctx) (primitive.ObjectID, error) {
	// Convert id parameter to objectId
	id, err := ConvertStringIdIntoObjectId(c.Params("id"))
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return id, nil
}

func RetrieveGetByIdRequestData(c *fiber.Ctx) (primitive.ObjectID, error) {
	return RetrieveDeleteRequestData(c)
}

func RetrieveResetPasswordRequestData(c *fiber.Ctx) (primitive.ObjectID, error) {
	return RetrieveGetByIdRequestData(c)
}

func IsRequestFromSameUser(c *fiber.Ctx) (bool, fiber.Error) {
	claims, err := security.ParseToken(ExtractToken(c))
	if err != nil {
		return false, fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	return claims.Id == c.Params("id"), fiber.Error{}
}

func RetrieveChangePasswordRequestData(c *fiber.Ctx) (primitive.ObjectID, models.ChangePasswordArgs, fiber.Error) {
	args := models.ChangePasswordArgs{}
	id, err := ConvertStringIdIntoObjectId(c.Params("id"))
	if err != nil {
		return primitive.ObjectID{}, args, fiber.Error{Code: fiber.StatusBadRequest, Message: err.Error()}
	}

	err2 := c.BodyParser(&args)
	if err2 != nil {
		return id, args, fiber.Error{Code: fiber.StatusBadRequest, Message: err.Error()}
	}

	return id, args, fiber.Error{}
}

func IsRequestFromAdmin(c *fiber.Ctx) (bool, error) {
	token := ExtractToken(c)

	claims, err := security.ParseToken(token)
	if err != nil {
		return false, err
	}

	return claims.IsAdmin, nil
}

func RetrieveIdFromToken(c *fiber.Ctx) (primitive.ObjectID, fiber.Error) {
	claims, err := security.ParseToken(ExtractToken(c))
	if err != nil {
		return primitive.ObjectID{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	id, err := ConvertStringIdIntoObjectId(claims.Id)
	if err != nil {
		return primitive.ObjectID{}, fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	return id, fiber.Error{}
}
