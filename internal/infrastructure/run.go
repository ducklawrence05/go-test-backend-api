package infrastructure

import (
	"log"
	"strconv"

	"github.com/ducklawrence05/go-test-backend-api/internal/infrastructure/initialization"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/ducklawrence05/go-test-backend-api/pkg/setting"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppContext struct {
	Config *setting.Config
	Logger *logger.LoggerZap
	Pgdb   *gorm.DB
}

func Run() {
	cfg, err := initialization.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	logger := initialization.InitLogger(cfg.Logger)
	logger.Info("Config log ok", zap.String("ok", "success"))

	pgDb := initialization.InitPostgres(cfg.Postgres, logger)

	appCtx := &AppContext{
		Config: cfg,
		Logger: logger,
		Pgdb:   pgDb,
	}

	r := initialization.InitRouter(appCtx.Config, appCtx.Logger, appCtx.Pgdb)
	r.Run(":" + strconv.Itoa(cfg.Server.Port))
}
