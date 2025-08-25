package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/str"
	"github.com/go-playground/validator/v10"
)

// Check field is not equal to another field
func NotEqualField(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	other := fl.Parent().FieldByName(fl.Param()).String()
	return field != other
}

var Validate = validator.New()

func init() {
	Validate.RegisterValidation("neqfield", NotEqualField)
}

func TranslateValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		msgs := []string{}
		for _, fe := range ve {
			var msg string
			fieldName := str.ToSnakeCase(fe.Field())
			paramName := str.ToSnakeCase(fe.Param())
			switch fe.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", fieldName)
			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters", fieldName, paramName)
			case "max":
				msg = fmt.Sprintf("%s must be at most %s characters", fieldName, paramName)
			case "eqfield":
				msg = fmt.Sprintf("%s must be equal to %s", fieldName, paramName)
			case "nefield":
				msg = fmt.Sprintf("%s must be different from %s", fieldName, paramName)
			case "len":
				msg = fmt.Sprintf("%s must be exactly %s characters long", fieldName, paramName)
			default:
				msg = fmt.Sprintf("%s: %s", fieldName, fe.Tag())
			}
			msgs = append(msgs, msg)
		}
		return strings.Join(msgs, "; ")
	}
	return err.Error()
}
