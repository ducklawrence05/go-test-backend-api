package app

import (
	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/internal/controller/http"
	"github.com/ducklawrence05/go-test-backend-api/internal/infrastructure/wire/user"
	"github.com/ducklawrence05/go-test-backend-api/internal/initialization"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
)

func Run(cfg *config.Config) {
	// logger
	logger := logger.New(cfg.Logger)
	logger.Info("Config log successfully")

	// postgres
	pgDb := initialization.InitPostgres(&cfg.Postgres, logger)
	logger.Info("Initializing Postgres successfully")

	// redis
	rdb := initialization.InitRedis(&cfg.Redis, logger)
	logger.Info("Initializing Redis successfully")

	// usecase
	userAuthManager := user.InitUserAuthManager(cfg, pgDb, rdb)
	userProfileManager := user.InitUserProfileManager(cfg, pgDb)

	// router
	routerCfg := &http.RouterConfig{
		Config: cfg,
		Db:     pgDb,
		Rdb:    rdb,
		Logger: logger,
	}

	router := http.InitRouter(routerCfg, userAuthManager, userProfileManager)

	server := initialization.NewServer(cfg.HTTP.Port, router)
	initialization.RunServer(server, logger)
}
