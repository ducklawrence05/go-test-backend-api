package errorcode

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// 400
	ErrInvalidEmail    = errors.New("this email is already exists")
	ErrInvalidUserName = errors.New("this username is already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidOTP      = errors.New("invalid otp")

	// 401
	ErrInvalidToken = errors.New("invalid token")

	// 404
	ErrUserNotFound = errors.New("user not found")
	ErrOTPNotFound  = errors.New("otp not found or expired")

	// 500
	ErrUnexpectedSigningToken = errors.New("unexpected signing token")
)

// Map code -> http code
var errorStatusMap = map[error]int{
	// 400
	ErrInvalidEmail:    http.StatusBadRequest,
	ErrInvalidUserName: http.StatusBadRequest,
	ErrInvalidPassword: http.StatusBadRequest,
	ErrInvalidOTP:      http.StatusBadRequest,

	// 401
	ErrInvalidToken: http.StatusUnauthorized,

	// 404
	ErrUserNotFound: http.StatusNotFound,
	ErrOTPNotFound:  http.StatusNotFound,

	// 500
	ErrUnexpectedSigningToken: http.StatusInternalServerError,
}

// utils write error
func JSONError(c *gin.Context, err error) {
	status, ok := errorStatusMap[err]
	if !ok {
		status = http.StatusInternalServerError
	}
	c.JSON(status, gin.H{
		"error": err.Error(),
	})
}
