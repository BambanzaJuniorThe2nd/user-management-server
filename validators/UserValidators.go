package validators

import (
	"server/models"
	"server/util"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateCreateByAdminArgs(args models.CreateByAdminArgs) error {
	err := validation.ValidateStruct(&args,
		// Name cannot be empty
		validation.Field(&args.Name, validation.Required.Error("name is required")),
		// Email cannot be empty, and must be a valid email
		validation.Field(&args.Email, validation.Required.Error("email is required"), is.Email.Error("Invalid email")),
		// Title cannot be empty
		validation.Field(&args.Title, validation.Required.Error("title is required")),
		// Birthdate cannot be empty, and must be a date string of the format "YYYY-MM-DD"
		validation.Field(&args.Birthdate, validation.Required.Error("birthdate is required"), validation.Date("2006-01-02").Error("Invalid date for birthdate. Format: YYYY-MM-DD")),
	)

	return util.ParseValidationError(err)
}
