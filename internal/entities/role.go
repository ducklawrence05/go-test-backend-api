package entities

type Role struct {
	ID          uint   `gorm:"column:id;type:int;primaryKey"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description;type:text"`
}

func (Role) TableName() string {
	return "roles"
}
