package refreshtoken

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/model"
	"github.com/google/uuid"
)

type Repository interface {
	GetByTokenAndUserID(token string, userID uuid.UUID) (*model.RefreshToken, error)
	Create(refreshToken *model.RefreshToken) error
	Revoke(userID uuid.UUID, token string) error
	DeleteByUserID(userID uuid.UUID) error
}
