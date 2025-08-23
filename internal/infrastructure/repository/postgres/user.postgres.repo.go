package postgres

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userPostgres struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repository.UserRepository {
	return &userPostgres{db: db}
}

func (pg *userPostgres) GetByID(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := pg.db.Preload("Role").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pg *userPostgres) GetByUsername(userName string) (*entities.User, error) {
	var user entities.User
	err := pg.db.Preload("Role").First(&user, "user_name = ?", userName).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pg *userPostgres) Create(user *entities.User) error {
	err := pg.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (pg *userPostgres) IsUserNameTaken(userName string, excludeUserID uuid.UUID) (bool, error) {
	var count int64
	err := pg.db.Model(&entities.User{}).
		Where("user_name = ? AND id != ?", userName, excludeUserID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (pg *userPostgres) Update(user *entities.User, fields map[string]any) error {
	err := pg.db.Model(&user).Updates(fields).Error
	if err != nil {
		return err
	}
	return nil
}

// hard delete
func (pg *userPostgres) DeleteByID(userID uuid.UUID) error {
	err := pg.db.Where("id = ?", userID).Delete(&entities.User{}).Error
	if err != nil {
		return err
	}
	return nil
}
