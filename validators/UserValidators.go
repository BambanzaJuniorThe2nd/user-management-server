package validators

import (
	"regexp"
	"server/models"
	"server/util"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// func ValidateCreateByAdminArgs(args models.CreateByAdminArgs) error {
// 	err := validation.ValidateStruct(&args,
// 		// Name cannot be empty
// 		validation.Field(&args.Name, validation.Required.Error("name is required")),
// 		// Email cannot be empty, and must be a valid email
// 		validation.Field(&args.Email, validation.Required.Error("email is required"), is.Email.Error("Invalid email")),
// 		// Title cannot be empty
// 		validation.Field(&args.Title, validation.Required.Error("title is required")),
// 		// Birthdate cannot be empty, and must be a date string of the format "YYYY-MM-DD"
// 		validation.Field(&args.Birthdate, validation.Required.Error("birthdate is required"), validation.Date("2006-01-02").Error("Invalid date for birthdate. Format: YYYY-MM-DD")),
// 	)

// 	return util.ParseValidationError(err)
// }

func ValidateLoginArgs(args models.LoginArgs) error {
	err := validation.ValidateStruct(&args,
		// Email cannot be empty, and must be a valid email
		validation.Field(&args.Email, validation.Required.Error("email is required"), is.Email.Error("Invalid email")),
		// Password cannot be empty
		validation.Field(&args.Password, validation.Required.Error("password is required")),
	)

	return util.ParseValidationError(err)
}

func ValidateCreateArgs(args models.CreateArgs) error {
	err := validation.ValidateStruct(&args,
		// Name cannot be empty
		validation.Field(&args.CreateByAdminArgs.Name, validation.Required.Error("name is required")),
		// Email cannot be empty, and must be a valid email
		validation.Field(&args.CreateByAdminArgs.Email, validation.Required.Error("email is required"), is.Email.Error("Invalid email")),
		// Title cannot be empty
		validation.Field(&args.CreateByAdminArgs.Title, validation.Required.Error("title is required")),
		// Birthdate cannot be empty, and must be a date string of the format "YYYY-MM-DD"
		validation.Field(&args.CreateByAdminArgs.Birthdate, validation.Required.Error("birthdate is required"), validation.Date("2006-01-02").Error("Invalid date for birthdate. Format: YYYY-MM-DD")),
		// IsAdmin must be set to false
		validation.Field(&args.CreateByAdminArgs.IsAdmin, validation.In(false).Error("isAdmin can only be set to false")),
		// Password cannot be empty
		validation.Field(
			&args.Password,
			validation.Required.Error("password is required"),
			validation.Match(regexp.MustCompile("[0-9]")).Error("password must contain at least one digit"),
			validation.Match(regexp.MustCompile("[a-z]")).Error("password must contain at least one lowercase letter"),
			validation.Match(regexp.MustCompile("[A-Z]")).Error("password must contain at least one uppercase letter"),
			validation.Match(regexp.MustCompile("[#?!@$%^&*-]")).Error("password must contain at least one special character"),
			validation.Match(regexp.MustCompile("[a-zA-Z0-9#?!@$%^&*-]{8,}$")).Error("password must have at least 8 characters")),
	)

	return util.ParseValidationError(err)
}

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
		// IsAdmin cannot be empty
		validation.Field(&args.IsAdmin),
	)

	return util.ParseValidationError(err)
}

func ValidateUpdateArgs(args models.UpdateArgs) error {
	err := validation.ValidateStruct(&args,
		// Name cannot be empty
		validation.Field(&args.CreateByAdminArgs.Name, validation.Required.Error("name is required")),
		// Email cannot be empty, and must be a valid email
		validation.Field(&args.CreateByAdminArgs.Email, validation.Required.Error("email is required"), is.Email.Error("Invalid email")),
		// Title cannot be empty
		validation.Field(&args.CreateByAdminArgs.Title, validation.Required.Error("title is required")),
		// Birthdate cannot be empty, and must be a date string of the format "YYYY-MM-DD"
		validation.Field(&args.CreateByAdminArgs.Birthdate, validation.Required.Error("birthdate is required"), validation.Date("2006-01-02").Error("Invalid date for birthdate. Format: YYYY-MM-DD")),
		// IsAdmin must be set to false
		validation.Field(&args.CreateByAdminArgs.IsAdmin, validation.In(false).Error("isAdmin can only be set to false")),
		// Password cannot be empty
		validation.Field(
			&args.Password,
			validation.Required.Error("password is required"),
			validation.Match(regexp.MustCompile("[0-9]")).Error("password must contain at least one digit"),
			validation.Match(regexp.MustCompile("[a-z]")).Error("password must contain at least one lowercase letter"),
			validation.Match(regexp.MustCompile("[A-Z]")).Error("password must contain at least one uppercase letter"),
			validation.Match(regexp.MustCompile("[#?!@$%^&*-]")).Error("password must contain at least one special character"),
			validation.Match(regexp.MustCompile("[a-zA-Z0-9#?!@$%^&*-]{8,}$")).Error("password must have at least 8 characters")),
	)

	return util.ParseValidationError(err)
}

func ValidateUpdateByAdminArgs(args models.UpdateByAdminArgs) error {
	err := validation.ValidateStruct(&args,
		// Name cannot be empty
		validation.Field(&args.Name, validation.Required.Error("name is required")),
		// Email cannot be empty, and must be a valid email
		validation.Field(&args.Email, validation.Required.Error("email is required"), is.Email.Error("Invalid email")),
		// Title cannot be empty
		validation.Field(&args.Title, validation.Required.Error("title is required")),
		// Birthdate cannot be empty, and must be a date string of the format "YYYY-MM-DD"
		validation.Field(&args.Birthdate, validation.Required.Error("birthdate is required"), validation.Date("2006-01-02").Error("Invalid date for birthdate. Format: YYYY-MM-DD")),
		// IsAdmin cannot be empty
		validation.Field(&args.IsAdmin),
	)

	return util.ParseValidationError(err)
}

func ValidateChangePasswordArgs(args models.ChangePasswordArgs) error {
	err := validation.ValidateStruct(&args,
		// Password cannot be empty
		validation.Field(
			&args.Password,
			validation.Required.Error("password is required"),
			validation.Match(regexp.MustCompile("[0-9]")).Error("password must contain at least one digit"),
			validation.Match(regexp.MustCompile("[a-z]")).Error("password must contain at least one lowercase letter"),
			validation.Match(regexp.MustCompile("[A-Z]")).Error("password must contain at least one uppercase letter"),
			validation.Match(regexp.MustCompile("[#?!@$%^&*-]")).Error("password must contain at least one special character"),
			validation.Match(regexp.MustCompile("[a-zA-Z0-9#?!@$%^&*-]{8,}$")).Error("password must have at least 8 characters")),
	)

	return util.ParseValidationError(err)
}

func ValidateCreateDefaultAdminArgs(args models.CreateDefaultAdminArgs) error {
	err := validation.ValidateStruct(&args,
		// Name cannot be empty
		validation.Field(&args.Name, validation.Required.Error("name is required")),
		// Email cannot be empty, and must be a valid email
		validation.Field(&args.Email, validation.Required.Error("email is required"), is.Email.Error("Invalid email")),
		// Title cannot be empty
		validation.Field(&args.Title, validation.Required.Error("title is required")),
		// Birthdate cannot be empty, and must be a date string of the format "YYYY-MM-DD"
		validation.Field(&args.Birthdate, validation.Required.Error("birthdate is required"), validation.Date("2006-01-02").Error("Invalid date for birthdate. Format: YYYY-MM-DD")),
		// IsAdmin must be set to false
		validation.Field(&args.IsAdmin, validation.In(false).Error("isAdmin can only be set to false")),
		// Password cannot be empty
		validation.Field(
			&args.Password,
			validation.Required.Error("password is required"),
			validation.Match(regexp.MustCompile("[0-9]")).Error("password must contain at least one digit"),
			validation.Match(regexp.MustCompile("[a-z]")).Error("password must contain at least one lowercase letter"),
			validation.Match(regexp.MustCompile("[A-Z]")).Error("password must contain at least one uppercase letter"),
			validation.Match(regexp.MustCompile("[#?!@$%^&*-]")).Error("password must contain at least one special character"),
			validation.Match(regexp.MustCompile("[a-zA-Z0-9#?!@$%^&*-]{8,}$")).Error("password must have at least 8 characters")),
	)

	return util.ParseValidationError(err)
}
