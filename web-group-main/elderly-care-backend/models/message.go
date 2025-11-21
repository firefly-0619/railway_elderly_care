package models

import "time"

type MessageType uint8

// MessageType 消息类型枚举
// @Description 消息内容的类型
// @Enum 0=text文本 1=image图片 2=file文件
const (
	Text MessageType = iota
	Image
	File
)

type Message struct {
	ID      uint      `gorm:"primarykey"`
	Time    time.Time `gorm:"index;"`
	From    uint      `gorm:"index:index_from_to"`
	To      uint      `gorm:"index:index_from_to"`
	Content string
	Type    MessageType
}

func (Message) TableName() string {
	return "message"
}
