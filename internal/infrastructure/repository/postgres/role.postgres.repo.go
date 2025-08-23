package postgres

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"gorm.io/gorm"
)

type rolePgRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) repository.RoleRepository {
	return &rolePgRepo{db: db}
}

func (r *rolePgRepo) GetByName(ctx context.Context, name string) (*entities.Role, error) {
	var role entities.Role
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
