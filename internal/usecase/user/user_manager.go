package user

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/google/uuid"
)

type (
	UserAuthManager interface {
		SendRegistrationOTP(ctx context.Context, email string) error
		VerifyRegistrationOTP(ctx context.Context, email, otp string) (string, error)
		CompleteRegistration(ctx context.Context, vo CreateUserVO) (string, string, error)
		Login(ctx context.Context, vo LoginUserVO) (string, string, error)
		Logout(ctx context.Context, vo LogoutUserVO) error
		RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	}

	UserProfileManager interface {
		GetMe(ctx context.Context, userID uuid.UUID) (*entities.User, error)
		UpdateMe(ctx context.Context, vo UpdateMeVO) (*entities.User, error)
		ChangePassword(ctx context.Context, vo ChangePasswordVO) error
		DeleteMe(ctx context.Context, userID uuid.UUID) error
	}
)
