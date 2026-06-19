package model

import (
	"time"

	"gorm.io/gorm"
)

// Tag represents a forum tag.
type Tag struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:20;uniqueIndex" json:"name"`
	Preset    bool           `gorm:"default:false" json:"preset"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName returns the table name for Tag.
func (Tag) TableName() string { return "tags" }
