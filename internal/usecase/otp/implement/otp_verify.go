package implement

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/constants/errorcode"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/otp"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/stringutils"
)

type otpVerifyManager struct {
	otpRepo repository.OTPRepository
}

func NewOTPVerifyManager(
	otpRepo repository.OTPRepository,
) otp.OTPVerifyManager {
	return &otpVerifyManager{
		otpRepo: otpRepo,
	}
}

// CanSendVerify implements otp.OTPVerifyManager.
func (o *otpVerifyManager) VerifyOTP(ctx context.Context, otp string, params otp.OTPParams) (bool, error) {
	hashedIndentifier := stringutils.HashString(params.Identifier, params.Secret)
	storeOtp, err := o.otpRepo.GetOTP(ctx, hashedIndentifier, params.OTPType)
	if err != nil {
		return false, err
	}

	if storeOtp != otp {
		// incr attempt
		attempts, err := o.otpRepo.IncrementAttempt(ctx, hashedIndentifier, params.OTPType, params.TTL)
		if err != nil {
			return false, err
		}

		if attempts >= int64(params.Limit) {
			o.otpRepo.DeleteOTP(ctx, hashedIndentifier, params.OTPType)
			o.otpRepo.ResetAttempt(ctx, hashedIndentifier, params.OTPType)
			return false, errorcode.ErrOTPTooManyAttempts
		}
		return false, errorcode.ErrInvalidOTP
	}

	o.otpRepo.DeleteOTP(ctx, hashedIndentifier, params.OTPType)
	o.otpRepo.ResetAttempt(ctx, hashedIndentifier, params.OTPType)
	return true, nil
}
