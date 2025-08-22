package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid" json:"user_id"`
	Token     string    `gorm:"column:token;type:text" json:"token"`
	IssuedAt  time.Time `gorm:"column:issued_at" json:"issued_at"`
	ExpiresAt time.Time `gorm:"column:expires_at" json:"expires_at"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	Revoked   bool      `gorm:"column:revoked" json:"revoked"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
