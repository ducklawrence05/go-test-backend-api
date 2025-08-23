package initialization

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/ducklawrence05/go-test-backend-api/pkg/setting"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres(pgCfg setting.PostgresSetting, logger *logger.LoggerZap) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		pgCfg.Host, pgCfg.Username, pgCfg.Password, pgCfg.Dbname, pgCfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("InitPostgres initialization error", zap.Error(err))
	}

	sqlDb, err := db.DB()
	if err != nil {
		logger.Fatal("Get sql.DB error", zap.Error(err))
	}
	if err := sqlDb.Ping(); err != nil {
		logger.Fatal("Database not reachable", zap.Error(err))
	}

	logger.Info("InitPostgres initialization success")

	SetPool(sqlDb, pgCfg)

	return db
}

func SetPool(sqlDb *sql.DB, pgCfg setting.PostgresSetting) {
	sqlDb.SetMaxIdleConns(pgCfg.MaxIdleConns)
	sqlDb.SetMaxOpenConns(pgCfg.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(pgCfg.ConnMaxLifetime * time.Second)
}
