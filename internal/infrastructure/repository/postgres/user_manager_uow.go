package postgres

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/uow"
	"gorm.io/gorm"
)

type userManagerPgUow struct {
	db *gorm.DB
}

func NewUserManagerUow(db *gorm.DB) uow.UserManagerUow {
	return &userManagerPgUow{db: db}
}

func (d *userManagerPgUow) Do(
	ctx context.Context,
	fn func(r uow.UserManagerRepoProvider) error,
) error {
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repoProvider := &repoProvider{tx: tx}
		err := fn(repoProvider)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
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
