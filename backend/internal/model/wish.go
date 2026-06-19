package model

import (
	"time"

	"gorm.io/gorm"
)

// Wish represents a wish list item.
type Wish struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:100;not null" json:"title"`
	Description string         `gorm:"size:500" json:"description"`
	Category    string         `gorm:"size:20;default:other" json:"category"`     // item / travel / experience / other
	Amount      *int64         `json:"amount"`                                     // nullable, in cents
	Priority    string         `gorm:"size:10;default:normal" json:"priority"`    // normal / important / urgent
	Type        string         `gorm:"size:10;default:personal" json:"type"`      // personal / family
	Status      string         `gorm:"size:10;default:pending" json:"status"`     // pending / agreed / achieved / abandoned
	CreatorID   uint           `gorm:"index;not null" json:"creator_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联
	Creator Member `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName returns the table name for Wish.
func (Wish) TableName() string { return "wishes" }

// WishVote represents a vote on a wish.
type WishVote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	WishID    uint      `gorm:"index;not null" json:"wish_id"`
	MemberID  uint      `gorm:"index;not null" json:"member_id"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	Member Member `gorm:"foreignKey:MemberID" json:"member,omitempty"`
}

// TableName returns the table name for WishVote.
func (WishVote) TableName() string { return "wish_votes" }
