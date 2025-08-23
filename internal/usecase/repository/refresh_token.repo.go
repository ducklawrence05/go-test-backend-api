package repository

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	GetByTokenAndUserID(token string, userID uuid.UUID) (*entities.RefreshToken, error)
	Create(refreshToken *entities.RefreshToken) error
	Revoke(userID uuid.UUID, token string) error
	DeleteByUserID(userID uuid.UUID) error
}
