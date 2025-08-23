package payload

import (
	"time"

	"github.com/google/uuid"
)

type RegisterUserPayLoad struct {
	UserName        string `json:"user_name" binding:"required"`
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Password        string `json:"password" binding:"required,min=8,max=30"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type LoginUserPayLoad struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=30"`
}

type LogoutUserPayLoad struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserInfoResponse struct {
	ID        uuid.UUID        `json:"id"`
	UserName  string           `json:"user_name"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	CreatedAt time.Time        `json:"created_at"`
	Role      RoleInfoResponse `json:"role"`
}

type RoleInfoResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateMeUserPayLoad struct {
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ChangePasswordPayLoad struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=30,neqfield=OldPassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

type RefreshTokenPayLoad struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
