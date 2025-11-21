package vo

import (
	"elderly-care-backend/models"
	"time"
)

type AccountVO struct {
	AccountType models.LoginType `gorm:"-" json:"account_type"`
	Avatar      string           `json:"avatar"`
	Nickname    string           `gorm:"size:25;Index" json:"nickname"`
	Sex         models.Sex       `json:"sex"`                              // 性别
	Phone       string           `gorm:"size:32;uniqueIndex" json:"phone"` // 手机号作为登录账号
	Password    string           `gorm:"size:64" json:"-"`
	Age         int              `json:"age"`
	Role        string           `gorm:"size:20;default:'user'" json:"role"` // 用户角色: user, volunteer, admin
	// 新增位置相关字段
	Latitude           float64    `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude          float64    `gorm:"type:decimal(11,8)" json:"longitude"`
	Address            string     `gorm:"size:255" json:"address"`
	LastLocationUpdate *time.Time `json:"last_location_update"`
}
