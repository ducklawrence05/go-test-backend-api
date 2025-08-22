package user

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/app"
	"github.com/ducklawrence05/go-test-backend-api/internal/middleware"
	"github.com/ducklawrence05/go-test-backend-api/internal/wire/user"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (*UserRouter) InitUserRouter(router *gin.RouterGroup, app *app.Application) {
	userHandler, _ := user.InitUserRouterHandler(app)
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
		[]byte(app.Config.JWT.AccessTokenKey), app.Logger))
	{
		userRouterPrivate.POST("/logout", userHandler.Logout)
		userRouterPrivate.GET("/me", userHandler.GetMe)
		userRouterPrivate.PATCH("/me", userHandler.UpdateMe)
		userRouterPrivate.PUT("/change-password", userHandler.ChangePassword)
		// hard delete
		userRouterPrivate.DELETE("/me", userHandler.DeleteMe)
	}
}
