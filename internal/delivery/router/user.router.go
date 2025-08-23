package router

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/delivery/middleware"
	"github.com/ducklawrence05/go-test-backend-api/internal/infrastructure/wire/user"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/ducklawrence05/go-test-backend-api/pkg/setting"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitUserRouter(router *gin.RouterGroup, config *setting.Config, logger *logger.LoggerZap, db *gorm.DB) {
	userHandler, _ := user.InitUserRouterHandler(config, db)
	// public
	userRouterPublic := router.Group("/user")
	{
		userRouterPublic.POST("/register", userHandler.Register)
		userRouterPublic.POST("/login", userHandler.Login)
		userRouterPublic.POST("/refresh-token", userHandler.RefreshToken)
	}
	// private
	userRouterPrivate := router.Group("/user")
	userRouterPrivate.Use(middleware.AccessTokenMiddleware(
		[]byte(config.JWT.AccessTokenKey), logger))
	{
		userRouterPrivate.POST("/logout", userHandler.Logout)
		userRouterPrivate.GET("/me", userHandler.GetMe)
		userRouterPrivate.PATCH("/me", userHandler.UpdateMe)
		userRouterPrivate.PUT("/change-password", userHandler.ChangePassword)
		// hard delete
		userRouterPrivate.DELETE("/me", userHandler.DeleteMe)
	}
}
