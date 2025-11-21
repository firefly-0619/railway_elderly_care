package vo

import "time"

type SOSRecordVO struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	TaskID      uint      `json:"task_id"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Address     string    `json:"address"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	Timeout     int       `json:"timeout,omitempty"`
}

type SOSResponseVO struct {
	SOSID   uint                     `json:"sos_id"`
	TaskID  uint                     `json:"task_id"`
	Matches []map[string]interface{} `json:"matches"`
	Timeout int                      `json:"timeout"`
}
