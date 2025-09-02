package entities

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid"`
	Token     string    `gorm:"column:token;type:text"`
	IssuedAt  time.Time `gorm:"column:issued_at"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Revoked   bool      `gorm:"column:revoked"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
