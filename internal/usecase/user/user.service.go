package user

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/uow"
	"github.com/ducklawrence05/go-test-backend-api/pkg/setting"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// interface
type Service interface {
	Register(ctx context.Context, vo CreateUserVO) (string, string, *utils.MyError)
	Login(ctx context.Context, vo LoginUserVO) (string, string, *utils.MyError)
	Logout(ctx context.Context, vo LogoutUserVO) *utils.MyError
	GetMe(ctx context.Context, userID uuid.UUID) (*entities.User, *utils.MyError)
	UpdateMe(ctx context.Context, vo UpdateMeVO) (*entities.User, *utils.MyError)
	ChangePassword(ctx context.Context, vo ChangePasswordVO) *utils.MyError
	DeleteMe(ctx context.Context, userID uuid.UUID) *utils.MyError
	RefreshToken(ctx context.Context, refreshToken string) (string, string, *utils.MyError)
}

// implement
type userService struct {
	config           *setting.Config
	uow              uow.UserServiceUow
	userRepo         repository.UserRepository
	roleRepo         repository.RoleRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewService(
	config *setting.Config,
	uow uow.UserServiceUow,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
) Service {
	return &userService{
		config:           config,
		uow:              uow,
		userRepo:         userRepo,
		roleRepo:         roleRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (us *userService) Register(ctx context.Context, vo CreateUserVO) (string, string, *utils.MyError) {
	var accessToken, refreshToken string

	myErr := us.uow.Do(ctx, func(r uow.UserServiceRepoProvider) *utils.MyError {
		var defaultRole *entities.Role

		g, ctx := errgroup.WithContext(ctx)

		// check if user exists
		g.Go(func() error {
			exists, err := r.UserRepository().IsUserNameTaken(ctx, vo.UserName, uuid.Nil)
			if err != nil {
				return err
			}
			if exists {
				return &utils.MyError{
					Msg:        "this username is already exists",
					StatusCode: http.StatusBadRequest,
				}
			}
			return nil
		})

		// get default role
		g.Go(func() error {
			var err error
			defaultRole, err = r.RoleRepository().GetByName(ctx, "user")
			if err != nil {
				return err
			}
			return nil
		})

		if myErr := utils.WaitErrGroup(g); myErr != nil {
			return myErr
		}

		// hash pass
		hashedPassword, err := utils.HashPassword(vo.Password)
		if err != nil {
			return &utils.MyError{
				Msg:        err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}

		// create use
		user := &entities.User{
			ID:        uuid.New(),
			UserName:  vo.UserName,
			FirstName: vo.FirstName,
			LastName:  vo.LastName,
			Password:  hashedPassword,
			IsActive:  true,
			CreatedAt: time.Now(),
			RoleID:    defaultRole.ID,
		}

		// gene ac and rt
		accessToken, refreshToken, err = utils.GenerateAcAndRtTokens(us.config, user.ID)
		if err != nil {
			return &utils.MyError{
				Msg:        err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}

		// decode rt to get exp and iat
		claims, err := utils.ValidateToken([]byte(us.config.JWT.RefreshTokenKey), refreshToken)
		if err != nil {
			return &utils.MyError{
				Msg:        err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}

		// insert user into db
		g.Go(func() error {
			err = r.UserRepository().Create(ctx, user)
			if err != nil {
				return err
			}
			return nil
		})

		// insert rt to into db
		g.Go(func() error {
			err = r.RefreshTokenRepository().Create(ctx, &entities.RefreshToken{
				ID:        uuid.New(),
				UserID:    user.ID,
				Token:     refreshToken,
				IssuedAt:  claims.IssuedAt.Time,
				ExpiresAt: claims.ExpiresAt.Time,
				CreatedAt: time.Now(),
				Revoked:   false,
			})
			if err != nil {
				return err
			}
			return nil
		})

		if myErr := utils.WaitErrGroup(g); myErr != nil {
			return myErr
		}

		// commit
		return nil
	})

	if myErr != nil {
		return "", "", myErr
	}
	return accessToken, refreshToken, nil
}

func (us *userService) Login(ctx context.Context, vo LoginUserVO) (string, string, *utils.MyError) {
	// get user from db
	user, err := us.userRepo.GetByUsername(ctx, vo.UserName)
	if err != nil || !utils.ComparePasswords(user.Password, []byte(vo.Password)) {
		return "", "", &utils.MyError{
			Msg:        "invalid username or password",
			StatusCode: http.StatusBadRequest,
		}
	}

	// gene ac and rt
	accessToken, refreshToken, err := utils.GenerateAcAndRtTokens(us.config, user.ID)
	if err != nil {
		return "", "", &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// decode rt to get exp and iat
	claims, err := utils.ValidateToken([]byte(us.config.JWT.RefreshTokenKey), refreshToken)
	if err != nil {
		return "", "", &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// insert rt to into db
	err = us.refreshTokenRepo.Create(ctx, &entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
		CreatedAt: time.Now(),
		Revoked:   false,
	})
	if err != nil {
		return "", "", &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return accessToken, refreshToken, nil
}

func (us *userService) Logout(ctx context.Context, vo LogoutUserVO) *utils.MyError {
	// decode rt
	claims, err := utils.ValidateToken([]byte(us.config.JWT.RefreshTokenKey), vo.RefreshToken)
	if err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}

	// compare userID from ac and rt
	if claims.UserID != vo.UserID {
		return &utils.MyError{
			Msg:        "refresh token is invalid",
			StatusCode: http.StatusUnauthorized,
		}
	}

	// check if revoked or not
	if _, err = us.refreshTokenRepo.GetByTokenAndUserID(ctx, vo.RefreshToken, vo.UserID); err != nil {
		return &utils.MyError{
			Msg:        "refresh token is invalid",
			StatusCode: http.StatusUnauthorized,
		}
	}

	// revoke
	err = us.refreshTokenRepo.Revoke(ctx, vo.RefreshToken, vo.UserID)
	if err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (us *userService) GetMe(ctx context.Context, userID uuid.UUID) (*entities.User, *utils.MyError) {
	user, err := us.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utils.MyError{
				Msg:        "user not found",
				StatusCode: http.StatusNotFound,
			}
		}
		return nil, &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return user, nil
}

func (us *userService) UpdateMe(ctx context.Context, vo UpdateMeVO) (*entities.User, *utils.MyError) {
	// get user
	user, err := us.userRepo.GetByID(ctx, vo.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utils.MyError{
				Msg:        "user not found",
				StatusCode: http.StatusNotFound,
			}
		}
		return nil, &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// check username unique
	if vo.UserName != "" {
		taken, err := us.userRepo.IsUserNameTaken(ctx, vo.UserName, user.ID)
		if err != nil {
			return nil, &utils.MyError{
				Msg:        err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}
		if taken {
			return nil, &utils.MyError{
				Msg:        "username already taken",
				StatusCode: http.StatusUnprocessableEntity,
			}
		}
		user.UserName = vo.UserName
	}

	// other fields
	if vo.FirstName != "" {
		user.FirstName = vo.FirstName
	}
	if vo.LastName != "" {
		user.LastName = vo.LastName
	}

	// update
	if err = us.userRepo.Update(ctx, user, map[string]any{
		"user_name":  user.UserName,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}); err != nil {
		return nil, &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return user, nil
}

func (us *userService) ChangePassword(ctx context.Context, vo ChangePasswordVO) *utils.MyError {
	// get user
	user, err := us.userRepo.GetByID(ctx, vo.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utils.MyError{
				Msg:        "user not found",
				StatusCode: http.StatusNotFound,
			}
		}
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// check old password
	if !utils.ComparePasswords(user.Password, []byte(vo.OldPassword)) {
		return &utils.MyError{
			Msg:        "invalid password",
			StatusCode: http.StatusBadRequest,
		}
	}

	// hash password
	hashedPassword, err := utils.HashPassword(vo.NewPassword)
	if err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// change password
	if err := us.userRepo.Update(ctx, user, map[string]any{
		"password": hashedPassword,
	}); err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (us *userService) DeleteMe(ctx context.Context, userID uuid.UUID) *utils.MyError {
	return us.uow.Do(ctx, func(r uow.UserServiceRepoProvider) *utils.MyError {
		// check user exists
		if _, err := r.UserRepository().GetByID(ctx, userID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &utils.MyError{
					Msg:        "user not found",
					StatusCode: http.StatusNotFound,
				}
			}
			return &utils.MyError{
				Msg:        err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}

		// hard delete all rt
		if err := r.RefreshTokenRepository().DeleteByUserID(ctx, userID); err != nil {
			return &utils.MyError{
				Msg:        err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}

		// hard delete user
		if err := r.UserRepository().DeleteByID(ctx, userID); err != nil {
			return &utils.MyError{
				Msg:        err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}

		return nil
	})
}

func (us *userService) RefreshToken(ctx context.Context, refreshToken string) (string, string, *utils.MyError) {
	var accessToken, newRefreshToken string
	myErr := us.uow.Do(ctx, func(r uow.UserServiceRepoProvider) *utils.MyError {
		// validate token
		claims, err := utils.ValidateToken([]byte(us.config.JWT.RefreshTokenKey), refreshToken)
		if err != nil {
			return &utils.MyError{
				Msg:        "refresh token is invalid",
				StatusCode: http.StatusUnauthorized,
			}
		}

		// check token in db
		if _, err := r.RefreshTokenRepository().GetByTokenAndUserID(
			ctx, refreshToken, claims.UserID,
		); err != nil {
			return &utils.MyError{
				Msg:        "refresh token is invalid",
				StatusCode: http.StatusUnauthorized,
			}
		}

		// gene ac and rt
		accessToken, newRefreshToken, err = utils.GenerateAcAndRtTokens(us.config, claims.UserID)
		if err != nil {
			return &utils.MyError{
				Msg:        err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}

		// decode rt to get exp and iat
		newClaims, err := utils.ValidateToken(
			[]byte(us.config.JWT.RefreshTokenKey),
			newRefreshToken,
		)
		if err != nil {
			return &utils.MyError{
				Msg:        err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}

		// insert rt to into db
		err = r.RefreshTokenRepository().Create(ctx, &entities.RefreshToken{
			ID:        uuid.New(),
			UserID:    newClaims.UserID,
			Token:     newRefreshToken,
			IssuedAt:  newClaims.IssuedAt.Time,
			ExpiresAt: newClaims.ExpiresAt.Time,
			CreatedAt: time.Now(),
			Revoked:   false,
		})
		if err != nil {
			return &utils.MyError{
				Msg:        err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}

		// revoke old rt
		err = r.RefreshTokenRepository().Revoke(ctx, refreshToken, newClaims.UserID)
		if err != nil {
			return &utils.MyError{
				Msg:        err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}

		return nil
	})

	if myErr != nil {
		return "", "", myErr
	}
	return accessToken, newRefreshToken, nil
}
