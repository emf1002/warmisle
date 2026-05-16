package model

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"size:20;uniqueIndex"`
	Preset    bool           `gorm:"default:false"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Tag) TableName() string { return "tags" }
