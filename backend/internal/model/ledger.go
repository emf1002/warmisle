package model

import (
	"time"

	"gorm.io/gorm"
)

type Ledger struct {
	ID         uint           `gorm:"primaryKey"`
	Amount     int64          `gorm:"not null"`
	Note       string         `gorm:"size:200"`
	CategoryID uint           `gorm:"index"`
	CreatorID  uint           `gorm:"index"`
	OccurredAt time.Time      `gorm:"index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	// 关联
	Category Category `gorm:"foreignKey:CategoryID"`
	Creator  Member   `gorm:"foreignKey:CreatorID"`
	Members  []Member `gorm:"many2many:ledger_members;"`
}

func (Ledger) TableName() string { return "ledgers" }

type LedgerMember struct {
	LedgerID uint `gorm:"primaryKey"`
	MemberID uint `gorm:"primaryKey"`
}

func (LedgerMember) TableName() string { return "ledger_members" }
