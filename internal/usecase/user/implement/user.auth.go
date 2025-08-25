package implement

import (
	"context"
	"errors"
	"time"

	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/internal/constants/errorcode"
	"github.com/ducklawrence05/go-test-backend-api/internal/constants/otptype"
	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/uow"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/user"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/jwt"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/otputils"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/password"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/sendto"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/str"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// implement
type userAuthManager struct {
	config           *config.Config
	logger           logger.Interface
	uow              uow.UserManagerUow
	otpRepo          repository.OTPRepository
	userRepo         repository.UserRepository
	roleRepo         repository.RoleRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewUserAuthManager(
	config *config.Config,
	logger logger.Interface,
	uow uow.UserManagerUow,
	otpRepo repository.OTPRepository,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
) user.UserAuthManager {
	return &userAuthManager{
		config:           config,
		uow:              uow,
		logger:           logger,
		otpRepo:          otpRepo,
		userRepo:         userRepo,
		roleRepo:         roleRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (m *userAuthManager) SendRegistrationOTP(ctx context.Context, email string) error {
	// check email exists
	exists, err := m.userRepo.IsEmailTaken(ctx, email, uuid.Nil)
	if err != nil {
		return err
	}
	if exists {
		return errorcode.ErrInvalidEmail
	}

	// gene otp
	otp := otputils.GenerateSecureOTP()
	// hash email
	hashedEmail := str.HashString(email, []byte(m.config.OTP.EmailVerifyKey))
	// save otp to redis
	if err := m.otpRepo.SetOTP(ctx, hashedEmail, otp, otptype.OTPEmailVerify,
		m.config.OTP.EmailVerifyExpiresIn); err != nil {
		return err
	}

	// send otp to email
	go func() {
		err := sendto.SendTemplateEmailOtp(&m.config.SMTP, []string{email},
			"register-otp-verification.html", map[string]any{"otp": otp})
		if err != nil {
			m.logger.Error("Send email error", zap.Error(err))
		}
	}()

	return nil
}

func (m *userAuthManager) VerifyRegistrationOTP(ctx context.Context, email, otp string) (string, error) {
	// hash email
	hashedEmail := str.HashString(email, []byte(m.config.OTP.EmailVerifyKey))

	// check otp in redis
	storedOtp, err := m.otpRepo.GetOTP(ctx, hashedEmail, otptype.OTPEmailVerify)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", errorcode.ErrOTPNotFound
		}
		return "", err
	}
	if storedOtp != otp {
		return "", errorcode.ErrInvalidOTP
	}

	// delete otp in redis
	if err := m.otpRepo.DeleteOTP(ctx, hashedEmail, otptype.OTPEmailVerify); err != nil {
		m.logger.Warn("Cannot delete old otp from Redis", zap.Error(err))
	}

	// gene jwt token
	token, err := jwt.GenerateEmailVerifiedToken(m.config, email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (m *userAuthManager) CompleteRegistration(ctx context.Context, vo user.CreateUserVO) (string, string, error) {
	g, gCtx := errgroup.WithContext(ctx)

	drChan := make(chan *entities.Role, 1)
	hpChan := make(chan string, 1)

	// check if username exists
	g.Go(func() error {
		exists, err := m.userRepo.IsUserNameTaken(gCtx, vo.UserName, uuid.Nil)
		if err != nil {
			return err
		}
		if exists {
			return errorcode.ErrInvalidUserName
		}
		return nil
	})

	// re-check if email exists
	g.Go(func() error {
		exists, err := m.userRepo.IsEmailTaken(gCtx, vo.Email, uuid.Nil)
		if err != nil {
			return err
		}
		if exists {
			return errorcode.ErrInvalidEmail
		}
		return nil
	})

	// get default role
	g.Go(func() error {
		dr, err := m.roleRepo.GetByName(gCtx, "user")
		if err != nil {
			return err
		}
		drChan <- dr
		return nil
	})

	// hash pass
	g.Go(func() error {
		hp, err := password.HashPassword(vo.Password)
		if err != nil {
			return err
		}
		hpChan <- hp
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", "", err
	}

	// create use
	userID := uuid.New()
	user := &entities.User{
		ID:        userID,
		Email:     vo.Email,
		UserName:  vo.UserName,
		FirstName: vo.FirstName,
		LastName:  vo.LastName,
		Password:  <-hpChan,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		DeletedAt: nil,

		RoleID: (<-drChan).ID,
	}

	var accessToken, refreshToken string
	// begin transaction
	err := m.uow.Do(ctx, func(r uow.UserManagerRepoProvider) error {
		// insert user into db
		err := r.UserRepository().Create(ctx, user)
		if err != nil {
			return err
		}

		// gene ac and rt
		accessToken, refreshToken, err = jwt.GenerateAcAndRtTokens(m.config, user.ID)
		if err != nil {
			return err
		}

		// decode rt to get exp and iat
		claims, err := jwt.ValidateToken([]byte(m.config.JWT.RefreshTokenKey),
			refreshToken, jwt.NewUserClaims)
		if err != nil {
			return err
		}

		if err := r.RefreshTokenRepository().Create(ctx, &entities.RefreshToken{
			ID:        uuid.New(),
			UserID:    user.ID,
			Token:     refreshToken,
			IssuedAt:  claims.IssuedAt.Time,
			ExpiresAt: claims.ExpiresAt.Time,
			CreatedAt: time.Now(),
			Revoked:   false,
		}); err != nil {
			return err
		}

		// commit
		return nil
	})

	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (m *userAuthManager) Login(ctx context.Context, vo user.LoginUserVO) (string, string, error) {
	// get user from db
	user, err := m.userRepo.GetByUsername(ctx, vo.EmailOrUsername)
	if err != nil || !password.ComparePasswords(user.Password, []byte(vo.Password)) {
		return "", "", errorcode.ErrInvalidPassword
	}

	// gene ac and rt
	accessToken, refreshToken, err := jwt.GenerateAcAndRtTokens(m.config, user.ID)
	if err != nil {
		return "", "", err
	}

	// decode rt to get exp and iat
	claims, err := jwt.ValidateToken([]byte(m.config.JWT.RefreshTokenKey),
		refreshToken, jwt.NewUserClaims)
	if err != nil {
		return "", "", err
	}

	// insert rt to into db
	err = m.refreshTokenRepo.Create(ctx, &entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
		CreatedAt: time.Now(),
		Revoked:   false,
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (m *userAuthManager) Logout(ctx context.Context, vo user.LogoutUserVO) error {
	// decode rt
	// claims, err := jwt.ValidateToken[jwt.UserClaims]([]byte(m.config.JWT.RefreshTokenKey), vo.RefreshToken)
	// if err != nil {
	// 	return err
	// }
	claims := &jwt.UserClaims{}

	// compare userID from ac and rt
	if claims.UserID != vo.UserID {
		return errorcode.ErrInvalidToken
	}

	// check if revoked or not
	if _, err := m.refreshTokenRepo.GetByTokenAndUserID(ctx, vo.RefreshToken, vo.UserID); err != nil {
		return errorcode.ErrInvalidToken
	}

	// revoke
	err := m.refreshTokenRepo.Revoke(ctx, vo.RefreshToken, vo.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (m *userAuthManager) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	var accessToken, newRefreshToken string
	err := m.uow.Do(ctx, func(r uow.UserManagerRepoProvider) error {
		// validate token
		// claims, err := jwt.ValidateToken[jwt.UserClaims]([]byte(m.config.JWT.RefreshTokenKey), refreshToken)
		// if err != nil {
		// 	return errorcode.ErrInvalidToken
		// }

		claims := &jwt.UserClaims{}

		// check token in db
		if _, err := r.RefreshTokenRepository().GetByTokenAndUserID(
			ctx, refreshToken, claims.UserID,
		); err != nil {
			return errorcode.ErrInvalidToken
		}

		// gene ac and rt
		var err error
		accessToken, newRefreshToken, err = jwt.GenerateAcAndRtTokens(m.config, claims.UserID)
		if err != nil {
			return err
		}

		// decode rt to get exp and iat
		// newClaims, err := jwt.ValidateToken[jwt.UserClaims]([]byte(m.config.JWT.RefreshTokenKey), newRefreshToken)
		// if err != nil {
		// 	return err
		// }
		newClaims := &jwt.UserClaims{}

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
			return err
		}

		// revoke old rt
		err = r.RefreshTokenRepository().Revoke(ctx, refreshToken, newClaims.UserID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", "", err
	}
	return accessToken, newRefreshToken, nil
}
