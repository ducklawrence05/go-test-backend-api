package repository

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(id uuid.UUID) (*entities.User, error)
	GetByUsername(userName string) (*entities.User, error)
	Create(user *entities.User) error
	Update(user *entities.User, fields map[string]any) error
	IsUserNameTaken(userName string, excludeUserID uuid.UUID) (bool, error)
	DeleteByID(userID uuid.UUID) error
}
