package uow

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
)

type UserManagerUow interface {
	Do(ctx context.Context, fn func(r UserManagerRepoProvider) error) error
}

type UserManagerRepoProvider interface {
	UserRepository() repository.UserRepository
	RefreshTokenRepository() repository.RefreshTokenRepository
}
