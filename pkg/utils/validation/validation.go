package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/stringutils"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func IsUserName(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	reAllowed := regexp.MustCompile(`^[a-zA-Z0-9._]{8,20}$`)
	if !reAllowed.MatchString(s) {
		return false
	}

	reRepeat := regexp.MustCompile(`[_.]{2}`)
	if reRepeat.MatchString(s) {
		return false
	}

	if s[0] == '.' || s[0] == '_' || s[len(s)-1] == '.' || s[len(s)-1] == '_' {
		return false
	}

	return true
}

func NotEqualField(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	other := fl.Parent().FieldByName(fl.Param()).String()
	return field != other
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// register for custom validation
		v.RegisterValidation("username", IsUserName)
		v.RegisterValidation("neqfield", NotEqualField)
	}
}

func TranslateValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		msgs := []string{}
		for _, fe := range ve {
			var msg string
			fieldName := stringutils.ToSnakeCase(fe.Field())
			paramName := stringutils.ToSnakeCase(fe.Param())
			switch fe.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", fieldName)
			case "email":
				msg = fmt.Sprintf("%s must be an email", fieldName)
			case "username":
				msg = fmt.Sprintf("%s must be 8–20 characters long and only contain letters, digits, dot (.), and underscore (_)", fieldName)
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
