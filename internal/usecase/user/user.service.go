package user

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils"
	"github.com/google/uuid"
)

type Service interface {
	Register(vo CreateUserVO) *utils.MyError
	Login(vo LoginUserVO) (string, string, *utils.MyError)
	Logout(vo LogoutUserVO) *utils.MyError
	GetMe(userID uuid.UUID) (*entities.User, *utils.MyError)
	UpdateMe(vo UpdateMeVO) (*entities.User, *utils.MyError)
	ChangePassword(vo ChangePasswordVO) *utils.MyError
	DeleteMe(userID uuid.UUID) *utils.MyError
	RefreshToken(refreshToken string) (string, string, *utils.MyError)
}
