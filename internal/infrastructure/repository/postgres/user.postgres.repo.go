package postgres

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userPgRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repository.UserRepository {
	return &userPgRepo{db: db}
}

func (r *userPgRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userPgRepo) GetByUsername(ctx context.Context, userName string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		First(&user, "user_name = ?", userName).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userPgRepo) Create(ctx context.Context, user *entities.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userPgRepo) IsUserNameTaken(ctx context.Context, userName string, excludeUserID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("user_name = ? AND id != ?", userName, excludeUserID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userPgRepo) Update(ctx context.Context, user *entities.User, fields map[string]any) error {
	err := r.db.WithContext(ctx).
		Model(&user).
		Updates(fields).Error
	if err != nil {
		return err
	}
	return nil
}

// hard delete
func (r *userPgRepo) DeleteByID(ctx context.Context, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", userID).
		Delete(&entities.User{}).Error
	if err != nil {
		return err
	}
	return nil
}
