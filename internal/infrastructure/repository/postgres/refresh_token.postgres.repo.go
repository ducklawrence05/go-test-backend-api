package postgres

import (
	"fmt"

	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type rtPostgres struct {
	db *gorm.DB
}

func NewRefreshTokenRepo(db *gorm.DB) repository.RefreshTokenRepository {
	return &rtPostgres{db: db}
}

func (pgr *rtPostgres) GetByTokenAndUserID(token string, userID uuid.UUID) (*entities.RefreshToken, error) {
	var refreshToken entities.RefreshToken
	err := pgr.db.Where("user_id = ? AND token = ? AND revoked = false", userID, token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (pgr *rtPostgres) Create(refreshToken *entities.RefreshToken) error {
	err := pgr.db.Create(refreshToken).Error
	if err != nil {
		return err
	}
	return nil
}

func (pgr *rtPostgres) Revoke(userID uuid.UUID, token string) error {
	result := pgr.db.Model(&entities.RefreshToken{}).
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

func (pgr *rtPostgres) DeleteByUserID(userID uuid.UUID) error {
	err := pgr.db.Where("user_id = ?", userID).Delete(&entities.RefreshToken{}).Error
	if err != nil {
		return err
	}
	return nil
}
