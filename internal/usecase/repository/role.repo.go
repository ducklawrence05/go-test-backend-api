package repository

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
)

type RoleRepository interface {
	GetByName(name string) (*entities.Role, error)
}
