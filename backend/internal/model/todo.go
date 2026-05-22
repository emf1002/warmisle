package model

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:100;not null" json:"title"`
	Description string         `gorm:"size:500" json:"description"`
	Priority    string         `gorm:"size:10;default:normal" json:"priority"`   // normal / important / urgent
	Status      string         `gorm:"size:10;default:pending" json:"status"`    // pending / completed
	AssigneeID  *uint          `gorm:"index" json:"assignee_id"`
	CreatorID   uint           `gorm:"index;not null" json:"creator_id"`
	DueDate     *time.Time     `json:"due_date"`
	CompletedAt *time.Time     `json:"completed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联
	Assignee *Member `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	Creator  Member  `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

func (Todo) TableName() string { return "todos" }

type TodoLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TodoID     uint      `gorm:"index;not null" json:"todo_id"`
	FieldName  string    `gorm:"size:50;not null" json:"field_name"`
	OldValue   string    `gorm:"size:500" json:"old_value"`
	NewValue   string    `gorm:"size:500" json:"new_value"`
	OperatorID uint      `gorm:"index;not null" json:"operator_id"`
	CreatedAt  time.Time `json:"created_at"`

	// 关联
	Operator Member `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
}

func (TodoLog) TableName() string { return "todo_logs" }
