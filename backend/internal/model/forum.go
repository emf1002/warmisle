package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primaryKey"`
	Content   string         `gorm:"size:1000;not null"`
	CreatorID uint           `gorm:"index;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// 关联
	Creator Member `gorm:"foreignKey:CreatorID"`
}

func (Post) TableName() string { return "posts" }

type Topic struct {
	ID        uint           `gorm:"primaryKey"`
	Title     string         `gorm:"size:100;not null"`
	Content   string         `gorm:"size:2000"`
	TagID     *uint          `gorm:"index"`
	CreatorID uint           `gorm:"index;not null"`
	IsPinned  bool           `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// 关联
	Tag     Tag    `gorm:"foreignKey:TagID"`
	Creator Member `gorm:"foreignKey:CreatorID"`
}

func (Topic) TableName() string { return "topics" }

type Vote struct {
	ID        uint           `gorm:"primaryKey"`
	Title     string         `gorm:"size:100;not null"`
	CreatorID uint           `gorm:"index;not null"`
	IsMulti   bool           `gorm:"default:false"`
	Deadline  *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// 关联
	Creator Member       `gorm:"foreignKey:CreatorID"`
	Options []VoteOption `gorm:"foreignKey:VoteID"`
}

func (Vote) TableName() string { return "votes" }

type VoteOption struct {
	ID        uint      `gorm:"primaryKey"`
	VoteID    uint      `gorm:"index;not null"`
	Content   string    `gorm:"size:50;not null"`
	SortOrder int       `gorm:"default:0"`
	CreatedAt time.Time
}

func (VoteOption) TableName() string { return "vote_options" }

type VoteRecord struct {
	ID        uint      `gorm:"primaryKey"`
	VoteID    uint      `gorm:"index;not null"`
	OptionID  uint      `gorm:"index;not null"`
	MemberID  uint      `gorm:"index;not null"`
	CreatedAt time.Time

	// 关联
	Member Member     `gorm:"foreignKey:MemberID"`
	Option VoteOption `gorm:"foreignKey:OptionID"`
}

func (VoteRecord) TableName() string { return "vote_records" }

type Comment struct {
	ID         uint           `gorm:"primaryKey"`
	TargetType string         `gorm:"size:10;not null"` // post / topic / wish
	TargetID   uint           `gorm:"index;not null"`
	ParentID   *uint          `gorm:"index"`
	Content    string         `gorm:"size:500;not null"`
	CreatorID  uint           `gorm:"index;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	// 关联
	Creator Member  `gorm:"foreignKey:CreatorID"`
	Parent  *Comment `gorm:"foreignKey:ParentID"`
}

func (Comment) TableName() string { return "comments" }

type Like struct {
	ID         uint      `gorm:"primaryKey"`
	TargetType string    `gorm:"size:10;not null"` // post / topic / comment
	TargetID   uint      `gorm:"index;not null"`
	MemberID   uint      `gorm:"index;not null"`
	CreatedAt  time.Time

	// 关联
	Member Member `gorm:"foreignKey:MemberID"`
}

func (Like) TableName() string { return "likes" }
