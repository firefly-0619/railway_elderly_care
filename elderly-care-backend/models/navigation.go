// models/navigation.go
package models

import "time"

// 导航会话
type NavigationSession struct {
	BaseModel
	HelperID    uint             `gorm:"not null" json:"helper_id"`   // 帮助者
	ElderlyID   uint             `gorm:"not null" json:"elderly_id"`  // 求助者
	Status      NavigationStatus `gorm:"not null" json:"status"`      // 会话状态
	RouteData   string           `gorm:"type:text" json:"route_data"` // 路线数据JSON
	StartedAt   time.Time        `json:"started_at"`                  // 开始时间
	CompletedAt *time.Time       `json:"completed_at"`                // 完成时间
}

type NavigationStatus int

const (
	NavigationActive    NavigationStatus = iota // 进行中
	NavigationCompleted                         // 已完成
	NavigationCancelled                         // 已取消
)

// 导航历史
type NavigationHistory struct {
	BaseModel
	SessionID uint    `gorm:"not null" json:"session_id"`
	UserID    uint    `gorm:"not null" json:"user_id"`
	FromLat   float64 `gorm:"type:decimal(10,8)" json:"from_lat"`
	FromLng   float64 `gorm:"type:decimal(11,8)" json:"from_lng"`
	ToLat     float64 `gorm:"type:decimal(10,8)" json:"to_lat"`
	ToLng     float64 `gorm:"type:decimal(11,8)" json:"to_lng"`
	Distance  int     `json:"distance"` // 总距离(米)
	Duration  int     `json:"duration"` // 总时长(秒)
}
