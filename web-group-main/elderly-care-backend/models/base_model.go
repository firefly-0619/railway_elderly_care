package models

import (
	"time"
)

type BaseModel struct {
	ID        uint `gorm:"primarykey;autoIncrement"` //自增ID
	CreatedAt time.Time
	UpdatedAt time.Time
}
