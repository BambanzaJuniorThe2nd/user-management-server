package validators

import (
	"server/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

func ValidateLoginArgs(args models.LoginArgs) error {
	err := validation.ValidateStruct(&args,
		// Email cannot be empty, and must be a valid email
		validation.Field(&args.Email, emailValidationRules...),
		// Password cannot be empty
		validation.Field(&args.Password, passwordValidationRules[0]),
	)

	return ParseValidationError(err)
}

func ValidateCreateByAdminArgs(args models.CreateByAdminArgs) error {
	err := validation.ValidateStruct(&args,
		// Name cannot be empty
		validation.Field(&args.Name, nameValidationRules...),
		// Email cannot be empty, and must be a valid email
		validation.Field(&args.Email, emailValidationRules...),
		// Title cannot be empty
		validation.Field(&args.Title, titleValidationRules...),
		// Birthdate cannot be empty, and must be a date string of the format "YYYY-MM-DD"
		validation.Field(&args.Birthdate, birthdateValidationRules...),
		// IsAdmin cannot be empty
		validation.Field(&args.IsAdmin),
	)

	return ParseValidationError(err)
}

func ValidateUpdateByAdminArgs(args models.UpdateByAdminArgs) error {
	err := validation.ValidateStruct(&args,
		// Name cannot be empty
		validation.Field(&args.Name, nameValidationRules...),
		// Email cannot be empty, and must be a valid email
		validation.Field(&args.Email, emailValidationRules...),
		// Title cannot be empty
		validation.Field(&args.Title, titleValidationRules...),
		// Birthdate cannot be empty, and must be a date string of the format "YYYY-MM-DD"
		validation.Field(&args.Birthdate, birthdateValidationRules...),
		// IsAdmin cannot be empty
		validation.Field(&args.IsAdmin),
	)

	return ParseValidationError(err)
}

func ValidateChangePasswordArgs(args models.ChangePasswordArgs) error {
	err := validation.ValidateStruct(&args,
		// Password cannot be empty
		validation.Field(&args.Password, passwordValidationRules...),
	)

	return ParseValidationError(err)
}

func ValidateCreateDefaultAdminArgs(args models.CreateDefaultAdminArgs) error {
	err := validation.ValidateStruct(&args,
		// Name cannot be empty
		validation.Field(&args.Name, nameValidationRules...),
		// Email cannot be empty, and must be a valid email
		validation.Field(&args.Email, emailValidationRules...),
		// Title cannot be empty
		validation.Field(&args.Title, titleValidationRules...),
		// Birthdate cannot be empty, and must be a date string of the format "YYYY-MM-DD"
		validation.Field(&args.Birthdate, birthdateValidationRules...),
		// IsAdmin must be set to false
		validation.Field(&args.IsAdmin),
		// Password cannot be empty
		validation.Field(&args.Password, passwordValidationRules...),
	)

	return ParseValidationError(err)
}
