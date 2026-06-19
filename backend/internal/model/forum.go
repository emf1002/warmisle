package model

import (
	"time"

	"gorm.io/gorm"
)

// Post represents a forum post.
type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Content   string         `gorm:"size:1000;not null" json:"content"`
	CreatorID uint           `gorm:"index;not null" json:"creator_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联
	Creator Member `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName returns the table name for Post.
func (Post) TableName() string { return "posts" }

// Topic represents a forum topic.
type Topic struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:100;not null" json:"title"`
	Content   string         `gorm:"size:2000" json:"content"`
	TagID     *uint          `gorm:"index" json:"tag_id"`
	CreatorID uint           `gorm:"index;not null" json:"creator_id"`
	IsPinned  bool           `gorm:"default:false" json:"is_pinned"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联
	Tag     Tag    `gorm:"foreignKey:TagID" json:"tag,omitempty"`
	Creator Member `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName returns the table name for Topic.
func (Topic) TableName() string { return "topics" }

// Vote represents a forum vote.
type Vote struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:100;not null" json:"title"`
	CreatorID uint           `gorm:"index;not null" json:"creator_id"`
	IsMulti   bool           `gorm:"default:false" json:"is_multi"`
	Deadline  *LocalTime     `json:"deadline"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联
	Creator Member       `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Options []VoteOption `gorm:"foreignKey:VoteID" json:"options,omitempty"`
}

// TableName returns the table name for Vote.
func (Vote) TableName() string { return "votes" }

// VoteOption represents a vote option.
type VoteOption struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VoteID    uint      `gorm:"index;not null" json:"vote_id"`
	Content   string    `gorm:"size:50;not null" json:"content"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName returns the table name for VoteOption.
func (VoteOption) TableName() string { return "vote_options" }

// VoteRecord represents a vote record.
type VoteRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VoteID    uint      `gorm:"index;not null" json:"vote_id"`
	OptionID  uint      `gorm:"index;not null" json:"option_id"`
	MemberID  uint      `gorm:"index;not null" json:"member_id"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	Member Member     `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	Option VoteOption `gorm:"foreignKey:OptionID" json:"option,omitempty"`
}

// TableName returns the table name for VoteRecord.
func (VoteRecord) TableName() string { return "vote_records" }

// Comment represents a forum comment.
type Comment struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	TargetType string         `gorm:"size:10;not null" json:"target_type"` // post / topic / wish
	TargetID   uint           `gorm:"index;not null" json:"target_id"`
	ParentID   *uint          `gorm:"index" json:"parent_id"`
	Content    string         `gorm:"size:500;not null" json:"content"`
	CreatorID  uint           `gorm:"index;not null" json:"creator_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联
	Creator Member   `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Parent  *Comment `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
}

// TableName returns the table name for Comment.
func (Comment) TableName() string { return "comments" }

// Like represents a like/dislike record.
type Like struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TargetType string    `gorm:"size:10;not null" json:"target_type"` // post / topic / comment
	TargetID   uint      `gorm:"index;not null" json:"target_id"`
	MemberID   uint      `gorm:"index;not null" json:"member_id"`
	CreatedAt  time.Time `json:"created_at"`

	// 关联
	Member Member `gorm:"foreignKey:MemberID" json:"member,omitempty"`
}

// TableName returns the table name for Like.
func (Like) TableName() string { return "likes" }
