package model

import (
	"time"

	"gorm.io/gorm"
)

type Wish struct {
	ID          uint           `gorm:"primaryKey"`
	Title       string         `gorm:"size:100;not null"`
	Description string         `gorm:"size:500"`
	Category    string         `gorm:"size:20;default:other"` // item / travel / experience / other
	Amount      *int64         // nullable, in cents
	Priority    string         `gorm:"size:10;default:normal"` // normal / important / urgent
	Type        string         `gorm:"size:10;default:personal"` // personal / family
	Status      string         `gorm:"size:10;default:pending"`  // pending / agreed / achieved / abandoned
	CreatorID   uint           `gorm:"index;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// 关联
	Creator Member `gorm:"foreignKey:CreatorID"`
}

func (Wish) TableName() string { return "wishes" }

type WishVote struct {
	ID        uint      `gorm:"primaryKey"`
	WishID    uint      `gorm:"index;not null"`
	MemberID  uint      `gorm:"index;not null"`
	CreatedAt time.Time

	// 关联
	Member Member `gorm:"foreignKey:MemberID"`
}

func (WishVote) TableName() string { return "wish_votes" }
