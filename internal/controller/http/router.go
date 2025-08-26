package http

import (
	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/internal/controller/http/v1/router"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/user"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RouterConfig struct {
	Config *config.Config
	Db     *gorm.DB
	Rdb    *redis.Client
	Logger logger.Interface
}

func InitRouter(
	routerCfg *RouterConfig,
	uRegistration user.UserRegistrationManager,
	uAuth user.UserAuthManager,
	uProfile user.UserProfileManager,
) *gin.Engine {
	var r *gin.Engine

	if routerCfg.Config.HTTP.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	MainGroup := r.Group("/v1/go-test")
	{
		router.InitUserRouter(MainGroup,
			&router.UserRouterConfig{
				Config: routerCfg.Config,
				Db:     routerCfg.Db,
				Rdb:    routerCfg.Rdb,
				Logger: routerCfg.Logger,
			},
			uRegistration, uAuth, uProfile,
		)
	}

	return r
}
