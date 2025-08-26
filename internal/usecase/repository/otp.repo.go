package repository

import (
	"context"
	"time"

	"github.com/ducklawrence05/go-test-backend-api/internal/constants/otptype"
)

type OTPRepository interface {
	SetOTP(ctx context.Context, identifier string, otp string,
		otpType otptype.OTPType, ttl time.Duration) error
	GetOTP(ctx context.Context, identifier string, otpType otptype.OTPType) (string, error)
	DeleteOTP(ctx context.Context, identifier string, otpType otptype.OTPType) error
	CountRateLimit(ctx context.Context, identifier string,
		otpType otptype.OTPType, ttl time.Duration) (int64, error)
	IncrementAttempt(ctx context.Context, identifier string,
		otptype otptype.OTPType, ttl time.Duration) (int64, error)
	ResetAttempt(ctx context.Context, identifier string, otptype otptype.OTPType) error
}
