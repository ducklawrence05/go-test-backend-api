package postgres

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"gorm.io/gorm"
)

type rolePostgres struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) repository.RoleRepository {
	return &rolePostgres{db: db}
}

func (pg *rolePostgres) GetByName(name string) (*entities.Role, error) {
	var role entities.Role
	err := pg.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
