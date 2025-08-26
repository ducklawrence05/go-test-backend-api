package implement

import (
	"context"
	"errors"
	"time"

	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/internal/constants/errorcode"
	"github.com/ducklawrence05/go-test-backend-api/internal/constants/jwtpurpose"
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
	"gorm.io/gorm"
)

type userRegistrationManager struct {
	config           *config.Config
	logger           logger.Interface
	uow              uow.UserManagerUow
	otpRepo          repository.OTPRepository
	userRepo         repository.UserRepository
	roleRepo         repository.RoleRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewUserRegistrationManager(
	config *config.Config,
	logger logger.Interface,
	uow uow.UserManagerUow,
	otpRepo repository.OTPRepository,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
) user.UserRegistrationManager {
	return &userRegistrationManager{
		config:           config,
		uow:              uow,
		logger:           logger,
		otpRepo:          otpRepo,
		userRepo:         userRepo,
		roleRepo:         roleRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (m *userRegistrationManager) SendRegistrationOTP(ctx context.Context, email string) error {
	// check email exists
	exists, err := m.userRepo.IsEmailTaken(ctx, email, uuid.Nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if exists {
		return errorcode.ErrInvalidEmail
	}

	// gene otp
	otp := otputils.GenerateSecureOTP()
	// hash email
	hashedEmail := str.HashString(email, []byte(m.config.OTP.RegisterKey))
	// save otp to redis
	if err := m.otpRepo.SetOTP(ctx, hashedEmail, otp, otptype.OTPRegister,
		m.config.OTP.RegisterExpiresIn); err != nil {
		return err
	}

	// send otp to email
	go func() {
		err := sendto.SendTemplateEmailOtp(&m.config.SMTP, []string{email},
			"otp-register-verify.html", map[string]any{"otp": otp})
		if err != nil {
			m.logger.Error("Send email error", zap.Error(err))
		}
	}()

	return nil
}

func (m *userRegistrationManager) VerifyRegistrationOTP(ctx context.Context, email, otp string) (string, error) {
	// hash email
	hashedEmail := str.HashString(email, []byte(m.config.OTP.RegisterKey))

	// check otp in redis
	storedOtp, err := m.otpRepo.GetOTP(ctx, hashedEmail, otptype.OTPRegister)
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
	if err := m.otpRepo.DeleteOTP(ctx, hashedEmail, otptype.OTPRegister); err != nil {
		m.logger.Warn("Cannot delete old otp from Redis", zap.Error(err))
	}

	// gene jwt token
	token, err := jwt.GenerateEmailToken(m.config, email,
		[]byte(m.config.JWT.RefreshTokenKey), m.config.JWT.RefreshTokenExpiresIn, jwtpurpose.JWTRegister)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (m *userRegistrationManager) CompleteRegistration(ctx context.Context, vo user.CreateUserVO) (string, string, error) {
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
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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
		IsActive:  true,

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
			refreshToken, jwtpurpose.JWTRefresh)
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
