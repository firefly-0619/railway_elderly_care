// routes/location_route.go
package routes

import (
	"elderly-care-backend/config"
	"elderly-care-backend/controllers"
	"elderly-care-backend/global"
	"elderly-care-backend/services"

	"github.com/gin-gonic/gin"
)

func LocationRoute(e *gin.Engine, wsService *services.WebSocketService) {
	// 从配置创建高德地图服务
	amapService := &services.AMapService{
		APIKey:  config.Config.Map.AMap.APIKey, // 从配置读取
		BaseURL: config.Config.Map.AMap.BaseURL,
	}

	// 创建位置服务，传入全局Db实例
	locationService := services.NewRealtimeLocationService(global.Db)

	// 创建控制器
	controller := controllers.NewLocationController(
		amapService,
		locationService,
		wsService,
	)

	locationRoute := e.Group("/location")
	{
		locationRoute.POST("/update", controller.UpdateLocation)
		locationRoute.GET("/user/:userId", controller.GetUserLocation)
		locationRoute.GET("/nearby", controller.GetNearbyUsers)
		locationRoute.GET("/navigation", controller.CalculateNavigation)
		locationRoute.GET("/navigation/to-target", controller.GetNavigationToTarget)
		locationRoute.GET("/navigation/user", controller.NavigateToUser)         // 新增
		locationRoute.GET("/navigation/location", controller.NavigateToLocation) // 新增
		locationRoute.GET("/history", controller.GetLocationHistory)             // 新增
		locationRoute.GET("/reverse-geocode", controller.ReverseGeocode)
	}

	// WebSocket路由
	e.GET("/ws", func(c *gin.Context) {
		wsService.HandleWebSocket(c.Writer, c.Request)
	})
}
