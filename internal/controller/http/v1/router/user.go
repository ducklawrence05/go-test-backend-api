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
	uProfileCtrl := controller.NewUserProfileController(mSet.ProfileManager)
	uRegistrationCtrl := controller.NewUserRegistrationController(mSet.RegistrationManager)
	uRestoreCtrl := controller.NewUserRestoreController(mSet.RestoreManager)
	uAuthCtrl := controller.NewUserAuthController(mSet.AuthManager)

	// ===== Public routes =====
	public := router.Group("/user")
	{
		public.POST("/login", uAuthCtrl.Login)
		public.POST("/refresh-token", uAuthCtrl.RefreshToken)
	}

	// Register route
	register := public.Group("/register")
	{
		register.POST("/send-email-otp", uRegistrationCtrl.SendRegistrationOTP)
		register.POST("/verify-email-otp", uRegistrationCtrl.VerifyRegistrationOTP)
		register.POST("/complete",
			middleware.ValidateToken([]byte(cfg.Config.JWT.RegisterTokenKey), jwtpurpose.Register, cfg.Logger),
			uRegistrationCtrl.Register,
		)
	}

	// Restore
	restore := public.Group("/restore")
	{
		restore.POST("/send-email-otp", uRestoreCtrl.SendRestoreOTP)
		restore.POST("/verify-email-otp", uRestoreCtrl.VerifyRestoreOTP)
		restore.POST("/complete",
			middleware.ValidateToken([]byte(cfg.Config.JWT.RestoreAccountTokenKey), jwtpurpose.Restore, cfg.Logger),
			uRestoreCtrl.Restore,
		)
	}

	// ===== Private routes (need access token) =====
	private := router.Group("/user")
	// middleware
	private.Use(middleware.ValidateToken([]byte(cfg.Config.JWT.AccessTokenKey), jwtpurpose.Access, cfg.Logger))
	// controller
	{
		private.POST("/logout", uAuthCtrl.Logout)
		private.GET("/me", uProfileCtrl.GetMe)
		private.PATCH("/me", uProfileCtrl.UpdateMe)
		private.PUT("/change-password", uProfileCtrl.ChangePassword)
		private.DELETE("/me", uProfileCtrl.DeleteMe)
	}
}
