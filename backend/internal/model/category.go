// Package model provides data models for the warmisle application.
package model

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a ledger category.
type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Type      string         `gorm:"size:10;not null" json:"type"` // income / expense
	Name      string         `gorm:"size:20;not null" json:"name"`
	Icon      string         `gorm:"size:30" json:"icon"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Preset    bool           `gorm:"default:false" json:"preset"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName returns the table name for Category.
func (Category) TableName() string { return "categories" }
