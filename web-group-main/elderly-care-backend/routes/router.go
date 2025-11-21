// routes/router.go
package routes

import (
	"elderly-care-backend/middlewares"
	"elderly-care-backend/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(middlewares.VerifyMiddleware())

	// 初始化WebSocket服务
	wsService := services.NewWebSocketService()
	go wsService.Start() // 启动WebSocket服务

	AccountRoute(r)
	ChatRoute(r)
	FileRoute(r)
	EvaluationRoute(r)
	TaskRoute(r)                // 新增
	SOSRoute(r)                 // 新增
	LocationRoute(r, wsService) // 新增定位路由

	return r
}
