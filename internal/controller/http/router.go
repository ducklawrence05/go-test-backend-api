package http

import (
	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/internal/controller/http/v1/router"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	Config *config.Config
	Logger logger.Interface
}

func NewRouter(
	routerCfg *RouterConfig,
	userManagerSet *router.UserManagerSet,
) *gin.Engine {
	r := gin.Default()

	MainGroup := r.Group("/v1")
	{
		router.NewUserRouter(
			MainGroup,
			&router.UserRouterConfig{
				Config: routerCfg.Config,
				Logger: routerCfg.Logger,
			},
			userManagerSet,
		)
	}

	return r
}
