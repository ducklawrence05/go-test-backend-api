package otp

import (
	"context"
)

type (
	OTPRateLimitManager interface {
		CanSendRateLimit(ctx context.Context, params OTPParams) (bool, error)
	}

	OTPVerifyManager interface {
		VerifyOTP(ctx context.Context, otp string, params OTPParams) (bool, error)
	}
)
