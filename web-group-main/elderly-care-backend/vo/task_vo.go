package vo

import "time"

type TaskVO struct {
	ID          uint      `json:"id"`
	CreatorID   uint      `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Status      string    `json:"status"`
	Reward      float64   `json:"reward"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Address     string    `json:"address"`
	Distance    float64   `json:"distance,omitempty"`
	IsEmergency bool      `json:"is_emergency,omitempty"`
	CreatorName string    `json:"creator_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskMatchVO struct {
	Task    TaskVO                   `json:"task"`
	Matches []map[string]interface{} `json:"matches"`
}

type NavigationVO struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
	Steps    []Step  `json:"steps"`
	StartLat float64 `json:"start_lat"`
	StartLng float64 `json:"start_lng"`
	EndLat   float64 `json:"end_lat"`
	EndLng   float64 `json:"end_lng"`
}

type Step struct {
	Instruction string  `json:"instruction"`
	Distance    float64 `json:"distance"`
	Road        string  `json:"road"`
}
