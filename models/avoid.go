package models

import (
	"time"

	"gorm.io/gorm"
)

type Avoid struct {
	gorm.Model
	ID            uint      `gorm:"not null" json:"id"`
	UserID        uint      `gorm:"not null" json:"user_id"`
	Name          string    `gorm:"not null" json:"name"`
	Duration      int       `gorm:"not null" json:"duration"`
	StartDate     time.Time `gorm:"not null" json:"start_date"`
	LastCheckedIn time.Time `json:"last_checked_in"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
}

type DailyCheck struct {
	gorm.Model
	AvoidID     uint      `gorm:"not null" json:"avoid_id"`
	CheckedDate time.Time `gorm:"not null" json:"checked_date"`
}
