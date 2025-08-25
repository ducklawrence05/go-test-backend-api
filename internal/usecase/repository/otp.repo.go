package repository

import (
	"context"
	"time"

	"github.com/ducklawrence05/go-test-backend-api/internal/constants/otptype"
)

type OTPRepository interface {
	SetOTP(ctx context.Context, identifier string, otp string, otpType otptype.OTPType, expiresIn time.Duration) error
	GetOTP(ctx context.Context, identifier string, otpType otptype.OTPType) (string, error)
	DeleteOTP(ctx context.Context, identifier string, otpType otptype.OTPType) error
}
