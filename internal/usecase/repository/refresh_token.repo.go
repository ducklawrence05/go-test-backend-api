package repository

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	GetByTokenAndUserID(ctx context.Context, token string, userID uuid.UUID) (*entities.RefreshToken, error)
	Create(ctx context.Context, refreshToken *entities.RefreshToken) error
	Revoke(ctx context.Context, token string, userID uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
