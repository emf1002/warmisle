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
	Members  []Member `gorm:"many2many:ledger_members;" json:"members,omitempty"`
}

func (Ledger) TableName() string { return "ledgers" }

type LedgerMember struct {
	LedgerID uint `gorm:"primaryKey" json:"ledger_id"`
	MemberID uint `gorm:"primaryKey" json:"member_id"`
}

func (LedgerMember) TableName() string { return "ledger_members" }
