package form

import (
	"time"
)

type User struct {
	ID       string     `gorm:"column:id;unique;not null"`
	Mobile   string     `gorm:"column:mobile;unique;not null"`
	Name     string     `gorm:"column:name;not null"`
	Birthday *time.Time `gorm:"column:birthday;default:''"`
	Gender   string     `gorm:"column:gender;not null"`
	Role     int        `gorm:"column:role;not null;default:'1'"`
	Address  string     `gorm:"column:address;default:''"`
	Salt     string     `gorm:"column:salt;not null"`
	Password string     `gorm:"column:salt;not null"`
}
