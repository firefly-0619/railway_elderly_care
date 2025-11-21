package music_vo

type PlaylistVO struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	AccountName string `json:"accountName"`
	AccountID   uint   `json:"accountID"`
	MusicCount  uint   `json:"musicCount"`
	ViewCount   uint   `json:"viewCount"`
}
