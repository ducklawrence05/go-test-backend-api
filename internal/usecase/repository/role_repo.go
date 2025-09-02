package repository

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
)

type RoleRepository interface {
	GetAll(ctx context.Context) ([]entities.Role, error)
	GetByName(ctx context.Context, name string) (*entities.Role, error)
}
