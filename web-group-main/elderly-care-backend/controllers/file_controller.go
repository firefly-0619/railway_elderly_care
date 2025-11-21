package controllers

import (
	"elderly-care-backend/common/constants"
	"elderly-care-backend/common/factories"
	"elderly-care-backend/vo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path"
)

type FileController struct {
}

// @Tags 文件模块
// @Summary 文件上传
// @Description 文件上传
// @Accept mpfd
// @Produce json
// @Security ApiKeyAuth
// @Param file formData file true "文件"
// @Success 200 {object} vo.ResponseVO "成功"
// @Failure 500 {object} vo.ResponseVO "失败"
// @Router /file/upload [post]
func (fc *FileController) UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Fail(constants.FILE_UPLOAD_ERROR))
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Fail(constants.FILE_OPEN_ERROR))
		return
	}
	client := factories.OssClientFactory.GetOssClient(factories.MINIO)
	url, err := client.Upload(factories.FILE_BUCKET, uuid.NewString()+path.Ext(fileHeader.Filename), file, fileHeader.Size)
	if err != nil {
		c.JSON(http.StatusBadGateway, vo.Fail(constants.UPLOAD_ERROR))
		return
	}
	c.JSON(http.StatusOK, vo.Success(url))
}
