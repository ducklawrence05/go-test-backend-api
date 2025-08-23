package entities

type Role struct {
	ID          uint   `gorm:"column:id;type:int;primaryKey" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description;type:text" json:"description"`
}

func (Role) TableName() string {
	return "roles"
}
