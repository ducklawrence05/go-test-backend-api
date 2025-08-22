package role

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/model"
)

type Repository interface {
	GetByName(name string) (*model.Role, error)
}
