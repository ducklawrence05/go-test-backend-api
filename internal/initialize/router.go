package initialize

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/app"
	"github.com/ducklawrence05/go-test-backend-api/internal/router"
	"github.com/gin-gonic/gin"
)

func InitRouter(app *app.Application) *gin.Engine {
	var r *gin.Engine

	if app.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	userRouter := router.RouterGroupApp.User

	MainGroup := r.Group("/v1/go-test")
	{
		userRouter.InitUserRouter(MainGroup, app)
	}

	return r
}
