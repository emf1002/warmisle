package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `gorm:"primaryKey"`
	Type      string         `gorm:"size:10;not null"` // income / expense
	Name      string         `gorm:"size:20;not null"`
	Icon      string         `gorm:"size:10"`
	SortOrder int            `gorm:"default:0"`
	Preset    bool           `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Category) TableName() string { return "categories" }
