package vo

import "time"

type ContactListVo struct {
	ContactID    uint      `json:"contact_id"`
	Nickname     string    `json:"nickname"`
	Avatar       string    `json:"avatar"`
	LastChatTime time.Time `json:"last_chat_time"`
}
