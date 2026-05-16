package model

import (
	"time"

	"gorm.io/gorm"
)

type Member struct {
	ID        uint           `gorm:"primaryKey"`
	Username  string         `gorm:"uniqueIndex;size:20"`
	Password  string         `gorm:"size:255"`
	Name      string         `gorm:"size:20"`
	Avatar    string         `gorm:"size:10;default:👨"`
	Role      string         `gorm:"size:10;default:member"`
	Status    string         `gorm:"size:10;default:active"`
	LastLogin *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Member) TableName() string { return "members" }
