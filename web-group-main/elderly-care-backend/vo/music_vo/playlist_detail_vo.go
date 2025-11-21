package music_vo

import "time"

type PlayListDetailVo struct {
	PlaylistVO
	Introduction string    `json:"introduction"`
	CreatedAt    time.Time `json:"createdAt"`
}
