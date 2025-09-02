package router

import (
	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"

	"github.com/ducklawrence05/go-test-backend-api/internal/constants/jwtpurpose"
	"github.com/ducklawrence05/go-test-backend-api/internal/controller/http/middleware"
	controller "github.com/ducklawrence05/go-test-backend-api/internal/controller/http/v1/controller/user"
	otpUC "github.com/ducklawrence05/go-test-backend-api/internal/usecase/otp"
	userUC "github.com/ducklawrence05/go-test-backend-api/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type UserRouterConfig struct {
	Config *config.Config
	Logger logger.Interface
}

type UserManagerSet struct {
	RegistrationManager userUC.UserRegistrationManager
	RestoreManager      userUC.UserRestoreManager
	AuthManager         userUC.UserAuthManager
	ProfileManager      userUC.UserProfileManager
	OTPRateLimitManager otpUC.OTPRateLimitManager
	OTPVerifyManager    otpUC.OTPVerifyManager
}

func NewUserRouter(
	router *gin.RouterGroup,
	cfg *UserRouterConfig,
	mSet *UserManagerSet,
) {
	// New controller
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
			middleware.ValidateToken(cfg.Logger, []byte(cfg.Config.JWT.RegisterTokenKey), jwtpurpose.Register),
			uRegistrationCtrl.Register,
		)
	}

	// Restore
	restore := public.Group("/restore")
	{
		restore.POST("/send-email-otp", uRestoreCtrl.SendRestoreOTP)
		restore.POST("/verify-email-otp", uRestoreCtrl.VerifyRestoreOTP)
		restore.POST("/complete",
			middleware.ValidateToken(cfg.Logger, []byte(cfg.Config.JWT.RestoreAccountTokenKey), jwtpurpose.Restore),
			uRestoreCtrl.Restore,
		)
	}

	// ===== Private routes (need access token) =====
	private := router.Group("/user")
	// middleware
	private.Use(middleware.ValidateToken(cfg.Logger, []byte(cfg.Config.JWT.AccessTokenKey), jwtpurpose.Access))
	// controller
	{
		private.POST("/logout", uAuthCtrl.Logout)
		private.GET("/me", uProfileCtrl.GetMe)
		private.PATCH("/me", uProfileCtrl.UpdateMe)
		private.PUT("/change-password", uProfileCtrl.ChangePassword)
		private.DELETE("/me", uProfileCtrl.DeleteMe)
	}
}
