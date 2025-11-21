package music_vo

type MusicListVo struct {
	ID       uint     `json:"id"`
	Title    string   `json:"title"`
	Author   []string `json:"author" gorm:"serializer:json"`
	AuthorId []uint   `json:"authorId" gorm:"serializer:json"`
	Cover    string   `json:"cover"`
	Duration string   `json:"duration"`
}
