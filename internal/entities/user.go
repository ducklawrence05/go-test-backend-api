package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserName  string    `gorm:"column:user_name;type:varchar(255)" json:"user_name"`
	FirstName string    `gorm:"column:first_name;type:varchar(255)" json:"first_name"`
	LastName  string    `gorm:"column:last_name;type:varchar(255)" json:"last_name"`
	Password  string    `gorm:"column:password;type:varchar(255)" json:"-"`
	IsActive  bool      `gorm:"column:is_active" json:"is_active"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	RoleID    uint      `gorm:"column:role_id;type:int;not null" json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleID" json:"role"`
}

func (User) TableName() string {
	return "users"
}
