package middleware

import (
	"net/http"
	"strings"

	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	jwtutils "github.com/ducklawrence05/go-test-backend-api/pkg/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func AccessTokenMiddleware[T jwt.Claims](secret []byte, logger logger.Interface, newClaims func() T) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Info("missing Authorization header")
			permissionDenied(c)
			return
		}

		// split bearer
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Warn("invalid Authorization header format")
			permissionDenied(c)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// validate token
		claims, err := jwtutils.ValidateToken(secret, token, newClaims)
		if err != nil {
			logger.Warn("failed to validate token", zap.Error(err))
			permissionDenied(c)
			return
		}

		switch v := any(claims).(type) {
		case *jwtutils.UserClaims:
			c.Set("userID", v.UserID)
		case *jwtutils.EmailClaims:
			c.Set("email", v.Email)
		default:
			logger.Warn("unknown claims type")
			permissionDenied(c)
			return
		}

		c.Next()
	}
}

func permissionDenied(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "permission denied",
	})
}
