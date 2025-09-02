//go:build wireinject

package role

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/infrastructure/repository/postgres"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/role"
	roleImpl "github.com/ducklawrence05/go-test-backend-api/internal/usecase/role/implement"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewRoleManager(db *gorm.DB) role.RoleManager {
	wire.Build(
		postgres.NewRoleRepo,
		roleImpl.NewRoleManager,
	)
	return nil
}
