package model

import (
	"time"

	"gorm.io/gorm"
)

type Member struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:20" json:"username"`
	Password  string         `gorm:"size:255" json:"-"`
	Name      string         `gorm:"size:20" json:"name"`
	Avatar    string         `gorm:"size:10;default:👨" json:"avatar"`
	Role      string         `gorm:"size:10;default:member" json:"role"`
	Status    string         `gorm:"size:10;default:active" json:"status"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Member) TableName() string { return "members" }
