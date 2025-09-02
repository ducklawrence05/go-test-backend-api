package initialization

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(pgCfg *config.Postgres, logger logger.Interface) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		pgCfg.Host, pgCfg.Username, pgCfg.Password, pgCfg.Dbname, pgCfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		logger.Fatal("NewPostgres initialization error", zap.Error(err))
	}

	sqlDb, err := db.DB()
	if err != nil {
		logger.Fatal("Get sql.DB error", zap.Error(err))
	}
	if err := sqlDb.Ping(); err != nil {
		logger.Fatal("Database not reachable", zap.Error(err))
	}

	SetPool(sqlDb, pgCfg)

	return db
}

func SetPool(sqlDb *sql.DB, pgCfg *config.Postgres) {
	sqlDb.SetMaxIdleConns(pgCfg.MaxIdleConns)
	sqlDb.SetMaxOpenConns(pgCfg.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(pgCfg.ConnMaxLifetime * time.Second)
}
