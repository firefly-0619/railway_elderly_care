package dto

type CreateTaskRequest struct {
	CreatorID   uint    `json:"creator_id" binding:"required"` // 改为uint
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Reward      float64 `json:"reward"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Address     string  `json:"address" binding:"required"`
	Deadline    *string `json:"deadline,omitempty"`
}

type AcceptTaskRequest struct {
	VolunteerID  uint    `json:"volunteer_id" binding:"required"` // 改为uint
	VolunteerLat float64 `json:"volunteer_lat" binding:"required"`
	VolunteerLng float64 `json:"volunteer_lng" binding:"required"`
}
