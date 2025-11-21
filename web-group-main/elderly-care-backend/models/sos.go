package models

import (
	"time"
)

type SOSRecord struct {
	BaseModel
	UserID      uint       `json:"user_id"`
	TaskID      uint       `json:"task_id"`
	Latitude    float64    `gorm:"type:decimal(10,8);not null" json:"latitude"`
	Longitude   float64    `gorm:"type:decimal(11,8);not null" json:"longitude"`
	Address     string     `gorm:"size:255;not null" json:"address"`
	Description string     `gorm:"type:text" json:"description"`
	Severity    string     `gorm:"size:20;default:'high'" json:"severity"`
	Status      string     `gorm:"size:50;default:'pending'" json:"status"`
	TimeoutAt   *time.Time `json:"timeout_at,omitempty"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`

	// 关联关系
	User Account `gorm:"foreignKey:UserID" json:"-"`
	Task Task    `gorm:"foreignKey:TaskID" json:"-"`
}

func (s *SOSRecord) TableName() string {
	return "sos_record"
}
