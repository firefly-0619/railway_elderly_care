package routes

import (
	"elderly-care-backend/controllers"
	"github.com/gin-gonic/gin"
)

func FileRoute(e *gin.Engine) {

	fileController := controllers.FileController{}
	fileGroup := e.Group("/file")
	{
		fileGroup.POST("/upload", fileController.UploadFile)
	}

}
