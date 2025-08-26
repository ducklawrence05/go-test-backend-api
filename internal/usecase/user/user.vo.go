package user

import "github.com/google/uuid"

type CreateUserVO struct {
	Email     string
	UserName  string
	FirstName string
	LastName  string
	Password  string
}

type RestoreUserVO struct {
	Email       string
	NewPassword string
}

type LoginUserVO struct {
	EmailOrUsername string
	Password        string
}

type LogoutUserVO struct {
	UserID       uuid.UUID
	RefreshToken string
}

type UpdateMeVO struct {
	UserID    uuid.UUID
	UserName  string
	FirstName string
	LastName  string
}

type ChangePasswordVO struct {
	UserID      uuid.UUID
	OldPassword string
	NewPassword string
}
