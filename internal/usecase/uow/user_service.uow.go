package uow

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils"
)

type UserServiceUow interface {
	Do(ctx context.Context, fn func(r UserServiceRepoProvider) *utils.MyError) *utils.MyError
}

type UserServiceRepoProvider interface {
	UserRepository() repository.UserRepository
	RoleRepository() repository.RoleRepository
	RefreshTokenRepository() repository.RefreshTokenRepository
}
