package router

import (
	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"

	"github.com/ducklawrence05/go-test-backend-api/internal/constants/jwtpurpose"
	"github.com/ducklawrence05/go-test-backend-api/internal/constants/otptype"
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
		register.POST("/send-email-otp",
			middleware.OTPSendRateLimit(cfg.Logger, mSet.OTPRateLimitManager,
				otpUC.OTPParams{
					Secret:  []byte(cfg.Config.OTP.RegisterKey),
					OTPType: otptype.Register,
					Limit:   cfg.Config.OTP.RegisterRateLimit,
					TTL:     cfg.Config.OTP.RegisterRateLimitTTL,
				},
			), uRegistrationCtrl.SendRegistrationOTP)
		register.POST("/verify-email-otp", uRegistrationCtrl.VerifyRegistrationOTP)
		register.POST("/complete",
			middleware.ValidateToken(cfg.Logger, []byte(cfg.Config.JWT.RegisterTokenKey), jwtpurpose.Register),
			uRegistrationCtrl.Register,
		)
	}

	// Restore
	restore := public.Group("/restore")
	{
		restore.POST("/send-email-otp",
			middleware.OTPSendRateLimit(cfg.Logger, mSet.OTPRateLimitManager,
				otpUC.OTPParams{
					Secret:  []byte(cfg.Config.OTP.RestoreAccountKey),
					OTPType: otptype.RestoreAccount,
					Limit:   cfg.Config.OTP.RestoreAccountRateLimit,
					TTL:     cfg.Config.OTP.RestoreAccountRateLimitTTL,
				},
			),
			uRestoreCtrl.SendRestoreOTP,
		)
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
