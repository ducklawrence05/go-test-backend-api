package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/str"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func NotEmail(fl validator.FieldLevel) bool {
	return !emailRegex.MatchString(fl.Field().String())
}

func NotEqualField(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	other := fl.Parent().FieldByName(fl.Param()).String()
	return field != other
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// đăng ký trên Gin's validator
		v.RegisterValidation("nemail", NotEmail)
		v.RegisterValidation("neqfield", NotEqualField)
	}
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
			case "email":
				msg = fmt.Sprintf("%s must be an email", fieldName)
			case "nemail":
				msg = fmt.Sprintf("%s must not be an email", fieldName)
			case "eqfield":
				msg = fmt.Sprintf("%s must be equal to %s", fieldName, paramName)
			case "nefield":
				msg = fmt.Sprintf("%s must be different from %s", fieldName, paramName)
			case "len":
				msg = fmt.Sprintf("%s must be exactly %s characters long", fieldName, paramName)
			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters", fieldName, paramName)
			case "max":
				msg = fmt.Sprintf("%s must be at most %s characters", fieldName, paramName)
			default:
				msg = fmt.Sprintf("%s: %s", fieldName, fe.Tag())
			}
			msgs = append(msgs, msg)
		}
		return strings.Join(msgs, "; ")
	}
	return err.Error()
}
