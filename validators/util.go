package validators

import (
	"errors"
	"regexp"
	"server/messages"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var nameValidationRules = []validation.Rule{
	validation.Required.Error(messages.ERROR_NAME_REQUIRED),
}

var emailValidationRules = []validation.Rule{
	validation.Required.Error(messages.ERROR_EMAIL_REQUIRED),
	is.Email.Error(messages.ERROR_INVALID_EMAIL),
}

var titleValidationRules = []validation.Rule{
	validation.Required.Error(messages.ERROR_TITLE_REQUIRED),
}

var birthdateValidationRules = []validation.Rule{
	validation.Required.Error(messages.ERROR_BIRTHDATE_REQUIRED), 
	validation.Date("2006-01-02").Error("Invalid date for birthdate. Format: YYYY-MM-DD"),
}

var isAdminValidationRules = []validation.Rule{
	validation.In(false).Error("isAdmin can only be set to false"),
}

var passwordValidationRules = []validation.Rule{
	validation.Required.Error(messages.ERROR_PASSWORD_REQUIRED),
	validation.Match(regexp.MustCompile("[0-9]")).Error("password must contain at least one digit"),
	validation.Match(regexp.MustCompile("[a-z]")).Error("password must contain at least one lowercase letter"),
	validation.Match(regexp.MustCompile("[A-Z]")).Error("password must contain at least one uppercase letter"),
	validation.Match(regexp.MustCompile("[#?!@$%^&*-]")).Error("password must contain at least one special character"),
	validation.Match(regexp.MustCompile("[a-zA-Z0-9#?!@$%^&*-]{8,}$")).Error("password must have at least 8 characters"),
}

func ParseValidationError(err error) error {
	if err != nil {
		errorList := strings.Split(err.Error(), ";")
		return errors.New(errorList[0])
	}
	return nil
}