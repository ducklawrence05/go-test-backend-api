package user

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgres struct {
	db *gorm.DB
}

func NewPostgres(db *gorm.DB) Repository {
	return &postgres{db: db}
}

func (pgr *postgres) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := pgr.db.Preload("Role").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pgr *postgres) GetByUsername(userName string) (*model.User, error) {
	var user model.User
	err := pgr.db.Preload("Role").First(&user, "user_name = ?", userName).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pgr *postgres) Create(user *model.User) error {
	err := pgr.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (pgr *postgres) IsUserNameTaken(userName string, excludeUserID uuid.UUID) (bool, error) {
	var count int64
	err := pgr.db.Model(&model.User{}).
		Where("user_name = ? AND id != ?", userName, excludeUserID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (pgr *postgres) Update(user *model.User, fields map[string]any) error {
	err := pgr.db.Model(&user).Updates(fields).Error
	if err != nil {
		return err
	}
	return nil
}

// hard delete
func (pgr *postgres) DeleteByID(userID uuid.UUID) error {
	err := pgr.db.Where("id = ?", userID).Delete(&model.User{}).Error
	if err != nil {
		return err
	}
	return nil
}
