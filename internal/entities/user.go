package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"column:id;type:uuid;primaryKey"`
	Email     string         `gorm:"column:email;type:varchar(255)"`
	UserName  string         `gorm:"column:user_name;type:varchar(255)"`
	FirstName string         `gorm:"column:first_name;type:varchar(255)"`
	LastName  string         `gorm:"column:last_name;type:varchar(255)"`
	Password  string         `gorm:"column:password;type:varchar(255)"`
	IsActive  bool           `gorm:"column:is_active"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	RoleID uint `gorm:"column:role_id;type:int"`
	Role   Role `gorm:"foreignKey:RoleID"`
}

func (User) TableName() string {
	return "users"
}
