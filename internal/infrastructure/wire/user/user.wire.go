//go:build wireinject

package user

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/delivery/handler"
	"github.com/ducklawrence05/go-test-backend-api/internal/infrastructure/repository/postgres"
	us "github.com/ducklawrence05/go-test-backend-api/internal/usecase/user"
	"github.com/ducklawrence05/go-test-backend-api/pkg/setting"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitUserRouterHandler(config *setting.Config, db *gorm.DB) (*handler.UserHandler, error) {
	wire.Build(
		postgres.NewUserRepo,
		postgres.NewRoleRepo,
		postgres.NewRefreshTokenRepo,
		us.NewService,
		handler.NewUserHandler,
	)
	return new(handler.UserHandler), nil
}
