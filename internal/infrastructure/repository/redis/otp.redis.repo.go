package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/ducklawrence05/go-test-backend-api/internal/constants/otptype"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/redis/go-redis/v9"
)

type otpRedisRepo struct {
	rdb *redis.Client
}

func NewOtpRepo(rdb *redis.Client) repository.OTPRepository {
	return &otpRedisRepo{rdb: rdb}
}

// SetOTP implements repository.OTPRepository.
func (o *otpRedisRepo) SetOTP(ctx context.Context, identifier string, otp string, otpType otptype.OTPType, expiresIn time.Duration) error {
	key := fmt.Sprintf("opt:%s:%s", identifier, otpType)
	return o.rdb.Set(ctx, key, otp, expiresIn).Err()
}

// GetOTP implements repository.OTPRepository.
func (o *otpRedisRepo) GetOTP(ctx context.Context, identifier string, otpType otptype.OTPType) (string, error) {
	key := fmt.Sprintf("opt:%s:%s", identifier, otpType)
	return o.rdb.Get(ctx, key).Result()
}

// DeleteOTP implements repository.OTPRepository.
func (o *otpRedisRepo) DeleteOTP(ctx context.Context, identifier string, otpType otptype.OTPType) error {
	key := fmt.Sprintf("opt:%s:%s", identifier, otpType)
	return o.rdb.Del(ctx, key).Err()
}
