package user

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/model"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(id uuid.UUID) (*model.User, error)
	GetByUsername(userName string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User, fields map[string]any) error
	IsUserNameTaken(userName string, excludeUserID uuid.UUID) (bool, error)
	DeleteByID(userID uuid.UUID) error
}
