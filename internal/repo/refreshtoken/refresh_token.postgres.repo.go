package refreshtoken

import (
	"fmt"

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

func (pgr *postgres) GetByTokenAndUserID(token string, userID uuid.UUID) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := pgr.db.Where("user_id = ? AND token = ? AND revoked = false", userID, token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (pgr *postgres) Create(refreshToken *model.RefreshToken) error {
	err := pgr.db.Create(refreshToken).Error
	if err != nil {
		return err
	}
	return nil
}

func (pgr *postgres) Revoke(userID uuid.UUID, token string) error {
	result := pgr.db.Model(&model.RefreshToken{}).
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

func (pgr *postgres) DeleteByUserID(userID uuid.UUID) error {
	err := pgr.db.Where("user_id = ?", userID).Delete(&model.RefreshToken{}).Error
	if err != nil {
		return err
	}
	return nil
}
