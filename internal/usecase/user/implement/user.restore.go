package implement

import (
	"context"
	"errors"

	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/internal/constants/errorcode"
	"github.com/ducklawrence05/go-test-backend-api/internal/constants/jwtpurpose"
	"github.com/ducklawrence05/go-test-backend-api/internal/constants/otptype"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/user"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/jwt"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/otputils"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/sendto"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/str"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type userRestoreManager struct {
	config           *config.Config
	logger           logger.Interface
	otpRepo          repository.OTPRepository
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewUserRestoreManager(
	config *config.Config,
	logger logger.Interface,
	otpRepo repository.OTPRepository,
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
) user.UserRestoreManager {
	return &userRestoreManager{
		config:           config,
		logger:           logger,
		otpRepo:          otpRepo,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

// SendRestoreOTP implements user.UserRestoreManager.
func (m *userRestoreManager) SendRestoreOTP(ctx context.Context, email string) error {
	// gene otp
	otp := otputils.GenerateSecureOTP()
	// hash email
	hashedEmail := str.HashString(email, []byte(m.config.OTP.RestoreAccountKey))
	// save otp to redis
	if err := m.otpRepo.SetOTP(ctx, hashedEmail, otp, otptype.OTPRestoreAccount,
		m.config.OTP.RestoreAccountExpiresIn); err != nil {
		return err
	}

	// send otp to email
	go func() {
		err := sendto.SendTemplateEmailOtp(&m.config.SMTP, []string{email},
			"otp-restore-account.html", map[string]any{"otp": otp})
		if err != nil {
			m.logger.Error("Send email error", zap.Error(err))
		}
	}()
	return nil
}

// VerifyRestoreOTP implements user.UserRestoreManager.
func (m *userRestoreManager) VerifyRestoreOTP(ctx context.Context, email string, otp string) (string, error) {
	// hash email
	hashedEmail := str.HashString(email, []byte(m.config.OTP.RestoreAccountKey))

	// check otp in redis
	storedOtp, err := m.otpRepo.GetOTP(ctx, hashedEmail, otptype.OTPRestoreAccount)
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
	if err := m.otpRepo.DeleteOTP(ctx, hashedEmail, otptype.OTPRestoreAccount); err != nil {
		m.logger.Warn("Cannot delete old otp from Redis", zap.Error(err))
	}

	// gene jwt token
	token, err := jwt.GenerateEmailToken(m.config, email,
		[]byte(m.config.JWT.RestoreAccountTokenKey), m.config.JWT.RestoreAccountTokenExpiresIn, jwtpurpose.JWTRestore)
	if err != nil {
		return "", err
	}

	return token, nil
}

// CompleteRestore implements user.UserRestoreManager.
func (m *userRestoreManager) CompleteRestore(ctx context.Context, vo user.RestoreUserVO) (string, string, error) {
	panic("unimplemented")
}
