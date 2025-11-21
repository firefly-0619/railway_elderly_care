package models

// 用户实时位置记录
type UserLocation struct {
	BaseModel
	UserID    uint    `gorm:"not null;index" json:"user_id"`
	Latitude  float64 `gorm:"type:decimal(10,8);not null" json:"latitude"`
	Longitude float64 `gorm:"type:decimal(11,8);not null" json:"longitude"`
	Address   string  `gorm:"size:255" json:"address"`
	Accuracy  float64 `json:"accuracy"` // 定位精度
	Speed     float64 `json:"speed"`    // 移动速度
	Heading   float64 `json:"heading"`  // 移动方向
}

func (ul *UserLocation) TableName() string {
	return "user_location"
}

// 实时位置更新请求
type LocationUpdate struct {
	UserID    uint    `json:"user_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
	Accuracy  float64 `json:"accuracy"`
	Speed     float64 `json:"speed"`
	Heading   float64 `json:"heading"`
}
