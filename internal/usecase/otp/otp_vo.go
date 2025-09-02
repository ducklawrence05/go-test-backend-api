package otp

import (
	"time"

	"github.com/ducklawrence05/go-test-backend-api/internal/constants/otptype"
)

type OTPParams struct {
	Identifier string
	Secret     []byte
	OTPType    otptype.OTPType
	Limit      int
	TTL        time.Duration
}
