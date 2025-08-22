package middleware

import (
	"net/http"
	"strings"

	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AccessTokenMiddleware(secret []byte, logger *logger.LoggerZap) gin.HandlerFunc {
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
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		// validate token
		claims, err := utils.ValidateToken(secret, accessToken)
		if err != nil {
			logger.Warn("failed to validate access token", zap.Error(err))
			permissionDenied(c)
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func permissionDenied(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "permission denied",
	})
}
