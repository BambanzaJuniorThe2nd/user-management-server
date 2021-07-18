package database

import "server/models"

var (
	ERROR_MESSAGE_SOMETHING_WENT_WRONG = "Something went wrong"
	ERROR_MESSAGE_USER_NOT_FOUND       = "User not found"
	ERROR_MESSAGE_EMAIL_ALREADY_IN_USE = "email already in use"
	ERROR_MESSAGE_LOGIN_FAILED         = "Login failed"
	ERROR_MESSAGE_ACCESS_RESTRICTED    = "Access limited to admins only"
	DATE_FORMAT                        = "2006-01-02"
)

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