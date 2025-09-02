package role

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
)

type RoleManager interface {
	GetAll(ctx context.Context) ([]entities.Role, error)
	GetByName(ctx context.Context, roleName string) (*entities.Role, error)
}
