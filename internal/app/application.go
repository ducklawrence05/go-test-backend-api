package app

import (
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/ducklawrence05/go-test-backend-api/pkg/setting"
	"gorm.io/gorm"
)

type Application struct {
	Config setting.Config
	Logger *logger.LoggerZap
	Pgdb   *gorm.DB
}
