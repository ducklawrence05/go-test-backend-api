//go:build wireinject

package otp

import (
	rdRepo "github.com/ducklawrence05/go-test-backend-api/internal/infrastructure/repository/redis"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/otp"
	otpImpl "github.com/ducklawrence05/go-test-backend-api/internal/usecase/otp/implement"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

func NewOTPRateLimitManager(rd *redis.Client) otp.OTPRateLimitManager {
	wire.Build(
		rdRepo.NewOtpRepo,
		otpImpl.NewOTPRateLimitManager,
	)
	return nil
}

func NewOTPVerifyManager(rd *redis.Client) otp.OTPVerifyManager {
	wire.Build(
		rdRepo.NewOtpRepo,
		otpImpl.NewOTPVerifyManager,
	)
	return nil
}
