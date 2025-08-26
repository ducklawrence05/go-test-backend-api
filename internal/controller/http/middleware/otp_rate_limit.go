package middleware

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/constants/errorcode"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/otp"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

func OTPSendRateLimit(
	logger logger.Interface, oManager otp.OTPRateLimitManager,
	params otp.OTPParams,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		if email == "" {
			errorcode.JSONError(c, errorcode.ErrInvalidEmail)
			c.Abort()
			return
		}

		ctx := c.Request.Context()

		params.Identifier = email
		canSend, err := oManager.CanSendRateLimit(ctx, params)
		if err != nil {
			errorcode.JSONError(c, err)
			c.Abort()
			return
		}

		if !canSend {
			errorcode.JSONError(c, errorcode.ErrOTPRateLimit)
			c.Abort()
			return
		}

		c.Next()
	}
}
