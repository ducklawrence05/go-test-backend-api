//go:build wireinject

package user

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/app"
	"github.com/ducklawrence05/go-test-backend-api/internal/handler"
	"github.com/ducklawrence05/go-test-backend-api/internal/repo/refreshtoken"
	"github.com/ducklawrence05/go-test-backend-api/internal/repo/role"
	"github.com/ducklawrence05/go-test-backend-api/internal/repo/user"
	us "github.com/ducklawrence05/go-test-backend-api/internal/service/user"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func ProvidePgDB(app *app.Application) *gorm.DB {
	return app.Pgdb
}

func InitUserRouterHandler(app *app.Application) (*handler.UserHandler, error) {
	wire.Build(
		ProvidePgDB,
		user.NewPostgres,
		role.NewPostgres,
		refreshtoken.NewPostgres,
		us.NewService,
		handler.NewUserHandler,
	)
	return new(handler.UserHandler), nil
}
