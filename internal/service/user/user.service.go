package user

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/payload"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils"
	"github.com/google/uuid"
)

type Service interface {
	Register(req payload.RegisterUserPayLoad) *utils.MyError
	Login(req payload.LoginUserPayLoad) (string, string, *utils.MyError)
	Logout(userID uuid.UUID, req payload.LogoutUserPayLoad) *utils.MyError
	GetMe(userID uuid.UUID) (*payload.GetMeUserResponse, *utils.MyError)
	UpdateMe(userID uuid.UUID, req payload.UpdateMeUserPayLoad) (*payload.GetMeUserResponse, *utils.MyError)
	ChangePassword(userID uuid.UUID, req payload.ChangePasswordPayLoad) *utils.MyError
	DeleteMe(userID uuid.UUID) *utils.MyError
	RefreshToken(req payload.RefreshTokenPayLoad) (string, string, *utils.MyError)
}
