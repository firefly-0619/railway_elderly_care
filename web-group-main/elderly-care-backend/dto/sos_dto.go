package dto

type SOSRequest struct {
	UserID      uint    `json:"user_id" binding:"required"` // 改为uint
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Address     string  `json:"address" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Severity    string  `json:"severity"`
}

type AcceptSOSRequest struct {
	VolunteerID uint `json:"volunteer_id" binding:"required"` // 改为uint
}

type ResolveSOSRequest struct {
	ResolvedBy uint `json:"resolved_by" binding:"required"` // 改为uint
}
