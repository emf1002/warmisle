package model

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          uint           `gorm:"primaryKey"`
	Title       string         `gorm:"size:100;not null"`
	Description string         `gorm:"size:500"`
	Priority    string         `gorm:"size:10;default:normal"` // normal / important / urgent
	Status      string         `gorm:"size:10;default:pending"` // pending / completed
	AssigneeID  *uint          `gorm:"index"`
	CreatorID   uint           `gorm:"index;not null"`
	DueDate     *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// 关联
	Assignee *Member `gorm:"foreignKey:AssigneeID"`
	Creator  Member  `gorm:"foreignKey:CreatorID"`
}

func (Todo) TableName() string { return "todos" }

type TodoLog struct {
	ID         uint      `gorm:"primaryKey"`
	TodoID     uint      `gorm:"index;not null"`
	FieldName  string    `gorm:"size:50;not null"`
	OldValue   string    `gorm:"size:500"`
	NewValue   string    `gorm:"size:500"`
	OperatorID uint      `gorm:"index;not null"`
	CreatedAt  time.Time

	// 关联
	Operator Member `gorm:"foreignKey:OperatorID"`
}

func (TodoLog) TableName() string { return "todo_logs" }
