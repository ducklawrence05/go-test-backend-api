package postgres

import (
	"context"
	"net/http"

	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/uow"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils"
	"gorm.io/gorm"
)

type userServicePgUow struct {
	db *gorm.DB
}

func NewUserServiceUow(db *gorm.DB) uow.UserServiceUow {
	return &userServicePgUow{db: db}
}

func (d *userServicePgUow) Do(ctx context.Context, fn func(r uow.UserServiceRepoProvider) *utils.MyError) *utils.MyError {
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repoProvider := &repoProvider{tx: tx}
		return fn(repoProvider)
	})

	if err != nil {
		if myErr, ok := err.(*utils.MyError); ok {
			return myErr
		}
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

type repoProvider struct {
	tx               *gorm.DB
	userRepo         repository.UserRepository
	roleRepo         repository.RoleRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func (r *repoProvider) UserRepository() repository.UserRepository {
	if r.userRepo == nil {
		r.userRepo = NewUserRepo(r.tx)
	}
	return r.userRepo
}

func (r *repoProvider) RoleRepository() repository.RoleRepository {
	if r.roleRepo == nil {
		r.roleRepo = NewRoleRepo(r.tx)
	}
	return r.roleRepo
}

func (r *repoProvider) RefreshTokenRepository() repository.RefreshTokenRepository {
	if r.refreshTokenRepo == nil {
		r.refreshTokenRepo = NewRefreshTokenRepo(r.tx)
	}
	return r.refreshTokenRepo
}
