package models

import "time"

type ContactList struct {
	BaseModel
	AccountID    uint      `gorm:"uniqueIndex:unique_account_id_contact_id"` //关联到某个账号
	ContactID    uint      `gorm:"uniqueIndex:unique_account_id_contact_id"` //联系人账号ID
	LastChatTime time.Time `gorm:"index"`                                    //上次聊天的时间
}

func (*ContactList) TableName() string {
	return "contact_list"
}
