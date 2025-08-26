package implement

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/role"
)

type roleManager struct {
	roleRepo repository.RoleRepository
}

func NewRoleManager(roleRepo repository.RoleRepository) role.RoleManager {
	return &roleManager{roleRepo: roleRepo}
}

// GetAll implements role.RoleManager.
func (r *roleManager) GetAll(ctx context.Context) ([]entities.Role, error) {
	roles, err := r.roleRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetByName implements role.RoleManager.
func (r *roleManager) GetByName(ctx context.Context, roleName string) (*entities.Role, error) {
	role, err := r.roleRepo.GetByName(ctx, roleName)
	if err != nil {
		return nil, err
	}
	return role, nil
}
