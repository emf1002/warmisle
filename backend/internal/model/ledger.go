package model

import (
	"time"

	"gorm.io/gorm"
)

type Ledger struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Amount     int64          `gorm:"not null" json:"amount"`
	Note       string         `gorm:"size:200" json:"note"`
	CategoryID uint           `gorm:"index" json:"category_id"`
	CreatorID  uint           `gorm:"index" json:"creator_id"`
	OccurredAt LocalTime      `gorm:"index" json:"occurred_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Creator  Member   `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

func (Ledger) TableName() string { return "ledgers" }
