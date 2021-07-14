package util

import (
	"errors"
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

func GetSafeUser(user models.User) models.User {
	return models.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Title:     user.Title,
		Birthdate: user.Birthdate,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
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

func ParseValidationError(err error) error {
	if err != nil {
		errorList := strings.Split(err.Error(), ";")
		return errors.New(errorList[0])
	}
	return nil
}

func RetrieveLoginRequestData(c *fiber.Ctx) (models.LoginArgs, error) {
	// Create an empty creds object
	// with the LoginArgs structure
	creds := models.LoginArgs{}

	err := c.BodyParser(&creds)
	return creds, err
}

func RetrieveCreateRequestData(c *fiber.Ctx, isAdmin bool) (interface{}, error) {
	if isAdmin {
		data := models.CreateByAdminArgs{}
		err := c.BodyParser(&data)
		return data, err
	} else {
		data := models.CreateArgs{}
		err := c.BodyParser(&data)
		return data, err
	}
}

func RetrieveUpdateRequestData(c *fiber.Ctx, isAdmin bool) (primitive.ObjectID, interface{}, error) {
	// Convert id parameter to objectId
	id, err := ConvertStringIdIntoObjectId(c.Params("id"))
	if err != nil {
		return primitive.ObjectID{}, models.UpdateArgs{}, err
	}

	if isAdmin {
		data := models.UpdateByAdminArgs{}
		err := c.BodyParser(&data)
		return id, data, err
	} else {
		data := models.UpdateArgs{}
		err := c.BodyParser(&data)
		return id, data, err
	}
}

func IsRequestFromAdmin(c *fiber.Ctx) (bool, error) {
	token := ExtractToken(c)

	claims, err := security.ParseToken(token)
	if err != nil {
		return false, err
	}

	return claims.IsAdmin, nil
}
