package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint    `gorm:"not null" json:"id"`
	Username string  `gorm:"unique;not null" json:"username"`
	Email    string  `gorm:"unique;not null" json:"email"`
	Password string  `gorm:"not null" json:"-"`
	Avoids   []Avoid `gorm:"foreignKey:UserID" json:"avoids"`
}
