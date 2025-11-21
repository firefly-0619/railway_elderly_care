package account_dto

import (
	"elderly-care-backend/models"
	"mime/multipart"
)

type AccountUpdateDTO struct {
	Nickname string               `form:"nickname"`
	Avatar   multipart.FileHeader `form:"avatar"`
	Sex      models.Sex           `form:"sex"`
	Age      int                  `form:"age"`
}
