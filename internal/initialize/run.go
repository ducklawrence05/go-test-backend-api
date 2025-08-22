package initialize

import (
	"log"
	"strconv"

	"github.com/ducklawrence05/go-test-backend-api/internal/app"
	"go.uber.org/zap"
)

func Run() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	logger := InitLogger(cfg.Logger)
	logger.Info("Config log ok", zap.String("ok", "success"))

	pgDb := InitPostgres(cfg.Postgres, logger)

	app := &app.Application{
		Config: *cfg,
		Logger: logger,
		Pgdb:   pgDb,
	}

	r := InitRouter(app)
	r.Run(":" + strconv.Itoa(cfg.Server.Port))
}
