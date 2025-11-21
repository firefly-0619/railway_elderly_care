package models

import (
	"time"
)

type Task struct {
	BaseModel
	CreatorID   uint       `json:"creator_id"`
	AssigneeID  *uint      `json:"assignee_id,omitempty"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	Description string     `gorm:"type:text;not null" json:"description"`
	Category    string     `gorm:"size:50;default:'other'" json:"category"`
	Status      string     `gorm:"size:50;default:'pending'" json:"status"`
	Reward      float64    `gorm:"type:decimal(10,2);default:0" json:"reward"`
	Latitude    float64    `gorm:"type:decimal(10,8);not null" json:"latitude"`
	Longitude   float64    `gorm:"type:decimal(11,8);not null" json:"longitude"`
	Address     string     `gorm:"size:255;not null" json:"address"`
	Deadline    *time.Time `json:"deadline,omitempty"`

	// 关联关系
	Creator  Account `gorm:"foreignKey:CreatorID" json:"-"`
	Assignee Account `gorm:"foreignKey:AssigneeID" json:"-"`
}

func (t *Task) TableName() string {
	return "task"
}
