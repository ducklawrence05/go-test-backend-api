package role

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/model"
	"gorm.io/gorm"
)

type postgres struct {
	db *gorm.DB
}

func NewPostgres(db *gorm.DB) Repository {
	return &postgres{db: db}
}

func (pgr *postgres) GetByName(name string) (*model.Role, error) {
	var role model.Role
	err := pgr.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
