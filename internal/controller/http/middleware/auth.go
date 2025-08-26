package middleware

import (
	"net/http"
	"strings"

	"github.com/ducklawrence05/go-test-backend-api/internal/constants/jwtpurpose"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	jwtutils "github.com/ducklawrence05/go-test-backend-api/pkg/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func ValidateToken(logger logger.Interface, secret []byte, purpose jwtpurpose.JWTPurpose) gin.HandlerFunc {
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
		claims, err := jwtutils.ValidateToken(secret, token, purpose)
		if err != nil {
			logger.Warn("failed to validate token", zap.Error(err))
			permissionDenied(c)
			return
		}

		switch claims.Purpose {
		case jwtpurpose.Access, jwtpurpose.Refresh:
			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "error when parsing claims subject to uuid",
				})
			}
			c.Set("userID", userID)
		case jwtpurpose.Register, jwtpurpose.Restore:
			c.Set("email", claims.Subject)
		}

		c.Next()
	}
}

func permissionDenied(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "permission denied",
	})
}
