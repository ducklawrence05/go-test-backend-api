package initialize

import (
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/ducklawrence05/go-test-backend-api/pkg/setting"
)

func InitLogger(loggerCfg setting.LoggerSetting) *logger.LoggerZap {
	return logger.NewLogger(loggerCfg)
}
