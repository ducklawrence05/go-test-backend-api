package implement

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/otp"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/stringutils"
)

type otpRateLimitManager struct {
	otpRepo repository.OTPRepository
}

func NewOTPRateLimitManager(
	otpRepo repository.OTPRepository,
) otp.OTPRateLimitManager {
	return &otpRateLimitManager{
		otpRepo: otpRepo,
	}
}

// CanSendRateLimit implements otp.OTPRateLimitManager.
func (o *otpRateLimitManager) CanSendRateLimit(ctx context.Context, params otp.OTPParams) (bool, error) {
	hashedIndentifier := stringutils.HashString(params.Identifier, params.Secret)
	count, err := o.otpRepo.CountRateLimit(ctx, hashedIndentifier, params.OTPType, params.TTL)
	if err != nil {
		return false, err
	}

	return count <= int64(params.Limit), nil
}
