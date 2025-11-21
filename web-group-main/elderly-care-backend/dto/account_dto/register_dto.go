package account_dto

import (
	"elderly-care-backend/models"
	"mime/multipart"
)

type RegisterDTO struct {
	Nickname string               `form:"nickname"`
	Phone    string               `form:"phone"` // 手机号作为登录账号
	Password string               `form:"password"`
	Avatar   multipart.FileHeader `form:"avatar"`
	Sex      models.Sex           `form:"sex"`
	Age      int                  `form:"age"`
	Role     string               `form:"role"`
}
