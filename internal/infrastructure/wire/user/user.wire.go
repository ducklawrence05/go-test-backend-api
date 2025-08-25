//go:build wireinject

package user

import (
	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/internal/infrastructure/repository/postgres"
	rdRepo "github.com/ducklawrence05/go-test-backend-api/internal/infrastructure/repository/redis"
	userInterface "github.com/ducklawrence05/go-test-backend-api/internal/usecase/user"
	userImpl "github.com/ducklawrence05/go-test-backend-api/internal/usecase/user/implement"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/redis/go-redis/v9"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitUserAuthManager(
	config *config.Config,
	db *gorm.DB,
	rdb *redis.Client,
	l logger.Interface,
) userInterface.UserAuthManager {
	wire.Build(
		rdRepo.NewOtpRepo,
		postgres.NewUserRepo,
		postgres.NewRoleRepo,
		postgres.NewRefreshTokenRepo,
		postgres.NewUserManagerUow,
		userImpl.NewUserAuthManager,
	)
	return nil
}

func InitUserProfileManager(
	config *config.Config,
	db *gorm.DB,
) userInterface.UserProfileManager {
	wire.Build(
		postgres.NewUserRepo,
		postgres.NewUserManagerUow,
		userImpl.NewUserProfileManager,
	)
	return nil
}
