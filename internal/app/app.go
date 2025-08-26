package app

import (
	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/internal/controller/http"
	"github.com/ducklawrence05/go-test-backend-api/internal/controller/http/v1/router"
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
	pgDb := initialization.InitPostgres(&cfg.Postgres, l)
	l.Info("Initializing Postgres successfully")

	// redis
	rdb := initialization.InitRedis(&cfg.Redis, l)
	l.Info("Initializing Redis successfully")

	// usecase
	userRegistrationManager := userWire.InitUserRegistrationManager(cfg, pgDb, rdb, l)
	userAuthManager := userWire.InitUserAuthManager(cfg, pgDb)
	userRestoreManager := userWire.InitUserRestoreManager(cfg, pgDb, rdb, l)
	userProfileManager := userWire.InitUserProfileManager(cfg, pgDb)
	roleManager := roleWire.InitRoleManager(pgDb)

	userManagerSet := &router.UserManagerSet{
		RegistrationManager: userRegistrationManager,
		RestoreManager:      userRestoreManager,
		AuthManager:         userAuthManager,
		ProfileManager:      userProfileManager,
	}

	// init role cache
	go initialization.InitRolesCache(roleManager, l)

	// router
	routerCfg := &http.RouterConfig{
		Config: cfg,
		Db:     pgDb,
		Rdb:    rdb,
		Logger: l,
	}

	router := http.InitRouter(routerCfg, userManagerSet)

	server := initialization.NewServer(cfg.HTTP.Port, router)
	initialization.RunServer(server, l)
}
