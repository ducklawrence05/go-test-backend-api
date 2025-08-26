package router

import (
	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"

	"github.com/ducklawrence05/go-test-backend-api/internal/constants/jwtpurpose"
	"github.com/ducklawrence05/go-test-backend-api/internal/controller/http/middleware"
	"github.com/ducklawrence05/go-test-backend-api/internal/controller/http/v1/controller"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRouterConfig struct {
	Config *config.Config
	Db     *gorm.DB
	Rdb    *redis.Client
	Logger logger.Interface
}

func InitUserRouter(
	router *gin.RouterGroup,
	urCfg *UserRouterConfig,
	registrationManager user.UserRegistrationManager,
	authManager user.UserAuthManager,
	profileManager user.UserProfileManager,
) {
	userController := controller.NewUserController(registrationManager, authManager, profileManager)

	// public
	userRouterPublic := router.Group("/user")
	{
		userRouterPublic.POST("/register",
			middleware.AccessTokenMiddleware(
				[]byte(urCfg.Config.JWT.RefreshTokenKey), urCfg.Logger, jwtpurpose.JWTRegister),
			userController.CompleteRegistration,
		)
		userRouterPublic.POST("/login", userController.Login)
		userRouterPublic.POST("/refresh-token", userController.RefreshToken)
	}

	// email
	userEmailRouter := userRouterPublic.Group("/email")
	{
		userEmailRouter.POST("/send-otp", userController.SendRegistrationOTP)
		userEmailRouter.POST("/verify-otp", userController.VerifyRegistrationOTP)
	}

	// private
	userRouterPrivate := router.Group("/user")
	userRouterPrivate.Use(middleware.AccessTokenMiddleware(
		[]byte(urCfg.Config.JWT.AccessTokenKey), urCfg.Logger, jwtpurpose.JWTAccess))
	{
		userRouterPrivate.POST("/logout", userController.Logout)
		userRouterPrivate.GET("/me", userController.GetMe)
		userRouterPrivate.PATCH("/me", userController.UpdateMe)
		userRouterPrivate.PUT("/change-password", userController.ChangePassword)
		// hard delete
		userRouterPrivate.DELETE("/me", userController.DeleteMe)
	}
}
