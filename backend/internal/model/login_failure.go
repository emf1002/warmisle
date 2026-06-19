package model

import "time"

// LoginFailure 登录失败记录（持久化到数据库，避免内存存储重启丢失）
type LoginFailure struct {
	Username    string     `gorm:"primaryKey;type:text"  json:"username"`
	FailedCount int        `gorm:"not null;default:0"    json:"failed_count"`
	LockedUntil *time.Time `gorm:"type:datetime"          json:"locked_until"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TableName returns the table name for LoginFailure.
func (LoginFailure) TableName() string { return "login_failures" }
