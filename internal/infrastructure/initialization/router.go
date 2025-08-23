package initialization

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/delivery/router"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/ducklawrence05/go-test-backend-api/pkg/setting"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(config *setting.Config, logger *logger.LoggerZap, db *gorm.DB) *gin.Engine {
	var r *gin.Engine

	if config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	MainGroup := r.Group("/v1/go-test")
	{
		router.InitUserRouter(MainGroup, config, logger, db)
	}

	return r
}
