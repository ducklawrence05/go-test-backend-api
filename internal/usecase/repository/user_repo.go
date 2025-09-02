package repository

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByUserNameOrEmail(ctx context.Context, identity string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User, fields map[string]any) error
	IsUserNameTaken(ctx context.Context, userName string, excludeUserID uuid.UUID) (bool, error)
	IsEmailTaken(ctx context.Context, email string, excludeUserID uuid.UUID) (bool, error)
	DeleteByID(ctx context.Context, userID uuid.UUID) error
}
