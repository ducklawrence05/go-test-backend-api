package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/ducklawrence05/go-test-backend-api/internal/constants/errorcode"
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
func (o *otpRedisRepo) SetOTP(ctx context.Context, identifier string, otp string, otpType otptype.OTPType, ttl time.Duration) error {
	key := fmt.Sprintf("otp:%s:%s", identifier, otpType)
	return o.rdb.Set(ctx, key, otp, ttl).Err()
}

// GetOTP implements repository.OTPRepository.
func (o *otpRedisRepo) GetOTP(ctx context.Context, identifier string, otpType otptype.OTPType) (string, error) {
	key := fmt.Sprintf("otp:%s:%s", identifier, otpType)
	otp, err := o.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errorcode.ErrOTPNotFound
		}
		return "", err
	}
	return otp, nil
}

// DeleteOTP implements repository.OTPRepository.
func (o *otpRedisRepo) DeleteOTP(ctx context.Context, identifier string, otpType otptype.OTPType) error {
	key := fmt.Sprintf("otp:%s:%s", identifier, otpType)
	return o.rdb.Del(ctx, key).Err()
}

// CountOTP implements repository.OTPRepository.
func (o *otpRedisRepo) CountRateLimit(ctx context.Context,
	identifier string, otpType otptype.OTPType, ttl time.Duration,
) (int64, error) {
	key := fmt.Sprintf("otp_rate_limit:%s:%s", identifier, otpType)
	count, err := o.rdb.Incr(ctx, key).Result()
	if err != nil {
		return -1, err
	}

	if count == 1 {
		o.rdb.Expire(ctx, key, ttl)
	}

	return count, nil
}

func (o *otpRedisRepo) IncrementAttempt(ctx context.Context, identifier string, otptype otptype.OTPType,
	ttl time.Duration) (int64, error) {
	key := fmt.Sprintf("otp_attempts:%s:%s", identifier, otptype)
	attemps, err := o.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if attemps == 1 {
		o.rdb.Expire(ctx, key, ttl)
	}
	return attemps, nil
}

func (o *otpRedisRepo) ResetAttempt(ctx context.Context, identifier string, otptype otptype.OTPType) error {
	key := fmt.Sprintf("otp_attempts:%s:%s", identifier, otptype)
	return o.rdb.Del(ctx, key).Err()
}
