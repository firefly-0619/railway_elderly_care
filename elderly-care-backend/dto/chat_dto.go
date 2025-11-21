package dto

import (
	"elderly-care-backend/models"
)

// MessageDTO 消息传输模型
// @Description 前端发送消息使用的数据结构
type MessageDTO struct {
	To      uint `gorm:"index:index_from_to"`
	Content string
	Type    models.MessageType
}
