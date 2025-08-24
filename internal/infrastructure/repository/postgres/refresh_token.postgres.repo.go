package postgres

import (
	"context"
	"fmt"

	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type refreshTokenPgRepo struct {
	db *gorm.DB
}

func NewRefreshTokenRepo(db *gorm.DB) repository.RefreshTokenRepository {
	return &refreshTokenPgRepo{db: db}
}

func (r *refreshTokenPgRepo) GetByTokenAndUserID(ctx context.Context, token string, userID uuid.UUID) (*entities.RefreshToken, error) {
	var refreshToken entities.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND token = ? AND revoked = false", userID, token).
		First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *refreshTokenPgRepo) Create(ctx context.Context, refreshToken *entities.RefreshToken) error {
	err := r.db.WithContext(ctx).Create(&refreshToken).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *refreshTokenPgRepo) Revoke(ctx context.Context, token string, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&entities.RefreshToken{}).
		Where("user_id = ? AND token = ? AND revoked = false", userID, token).
		Update("revoked", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("refresh token not found or already revoked")
	}

	return nil
}

func (r *refreshTokenPgRepo) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&entities.RefreshToken{}).Error
	if err != nil {
		return err
	}
	return nil
}
