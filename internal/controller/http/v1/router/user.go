package router

import (
	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"

	"github.com/ducklawrence05/go-test-backend-api/internal/constants/jwtpurpose"
	"github.com/ducklawrence05/go-test-backend-api/internal/controller/http/middleware"
	controller "github.com/ducklawrence05/go-test-backend-api/internal/controller/http/v1/controller/user"
	usecase "github.com/ducklawrence05/go-test-backend-api/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type UserRouterConfig struct {
	Config *config.Config
	Logger logger.Interface
}

type UserManagerSet struct {
	RegistrationManager usecase.UserRegistrationManager
	RestoreManager      usecase.UserRestoreManager
	AuthManager         usecase.UserAuthManager
	ProfileManager      usecase.UserProfileManager
}

func InitUserRouter(
	router *gin.RouterGroup,
	cfg *UserRouterConfig,
	mSet *UserManagerSet,
) {
	// Init controller
	uProfileController := controller.NewUserProfileController(mSet.ProfileManager)
	uRegistrationController := controller.NewUserRegistrationController(mSet.RegistrationManager)
	uRestoreController := controller.NewUserRestoreController(mSet.RestoreManager)
	uAuthController := controller.NewUserAuthController(mSet.AuthManager)

	userRouter := router.Group("/user")

	// --- Public routes ---
	userRouter.POST("/login", uAuthController.Login)
	userRouter.POST("/refresh-token", uAuthController.RefreshToken)

	// Register route
	registerRouter := userRouter.Group("/register")
	{
		registerRouter.POST("/send-email-otp", uRegistrationController.SendRegistrationOTP)
		registerRouter.POST("/verify-email-otp", uRegistrationController.VerifyRegistrationOTP)
		registerRouter.POST("/completion",
			middleware.ValidateToken(
				[]byte(cfg.Config.JWT.RegisterTokenKey), cfg.Logger, jwtpurpose.Register),
			uRegistrationController.Register,
		)
	}

	// Restore
	restoreRouter := userRouter.Group("/restore")
	{
		restoreRouter.POST("/send-email-otp", uRestoreController.SendRestoreOTP)
		restoreRouter.POST("/verify-email-otp", uRestoreController.VerifyRestoreOTP)
		restoreRouter.POST("/completion",
			middleware.ValidateToken(
				[]byte(cfg.Config.JWT.RestoreAccountTokenKey), cfg.Logger, jwtpurpose.Restore),
			uRestoreController.Restore,
		)
	}

	// --- Private routes (need access token) ---
	userRouter.Use(
		middleware.ValidateToken(
			[]byte(cfg.Config.JWT.AccessTokenKey), cfg.Logger, jwtpurpose.Access),
	)

	userRouter.POST("/logout", uAuthController.Logout)
	userRouter.GET("/me", uProfileController.GetMe)
	userRouter.PATCH("/me", uProfileController.UpdateMe)
	userRouter.PUT("/change-password", uProfileController.ChangePassword)
	userRouter.DELETE("/me", uProfileController.DeleteMe)
}
