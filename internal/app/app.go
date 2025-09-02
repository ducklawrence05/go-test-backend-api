package app

import (
	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/internal/controller/http"
	"github.com/ducklawrence05/go-test-backend-api/internal/controller/http/v1/router"
	otpWire "github.com/ducklawrence05/go-test-backend-api/internal/infrastructure/wire/otp"
	roleWire "github.com/ducklawrence05/go-test-backend-api/internal/infrastructure/wire/role"
	userWire "github.com/ducklawrence05/go-test-backend-api/internal/infrastructure/wire/user"
	"github.com/ducklawrence05/go-test-backend-api/internal/initialization"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
)

func Run(cfg *config.Config) {
	// logger
	l := logger.New(cfg.Logger)
	l.Info("Config log successfully")

	// postgres
	pgDb := initialization.NewPostgres(&cfg.Postgres, l)
	l.Info("Init Postgres successfully")

	// redis
	rdb := initialization.NewRedis(&cfg.Redis, l)
	l.Info("Init Redis successfully")

	// ===== usecase =====
	// user
	userRegistrationManager := userWire.NewUserRegistrationManager(cfg, pgDb, rdb, l)
	userAuthManager := userWire.NewUserAuthManager(cfg, pgDb)
	userRestoreManager := userWire.NewUserRestoreManager(cfg, pgDb, rdb, l)
	userProfileManager := userWire.NewUserProfileManager(cfg, pgDb)
	// role
	roleManager := roleWire.NewRoleManager(pgDb)
	// otp
	otpRateLimitManager := otpWire.NewOTPRateLimitManager(rdb)
	otpVerifyManager := otpWire.NewOTPVerifyManager(rdb)

	// init role cache
	go initialization.NewRolesCache(roleManager, l)

	// ===== router =====
	routerCfg := &http.RouterConfig{
		Config: cfg,
		Logger: l,
	}

	userManagerSet := &router.UserManagerSet{
		RegistrationManager: userRegistrationManager,
		RestoreManager:      userRestoreManager,
		AuthManager:         userAuthManager,
		ProfileManager:      userProfileManager,
		OTPRateLimitManager: otpRateLimitManager,
		OTPVerifyManager:    otpVerifyManager,
	}

	router := http.NewRouter(routerCfg, userManagerSet)

	server := initialization.NewServer(cfg.HTTP.Port, router)
	initialization.RunServer(server, l)
}
