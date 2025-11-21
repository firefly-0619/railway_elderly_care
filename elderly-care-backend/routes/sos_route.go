package routes

import (
	"elderly-care-backend/controllers"

	"github.com/gin-gonic/gin"
)

func SOSRoute(e *gin.Engine) {
	controller := controllers.NewSOSController()
	sosRoute := e.Group("/sos")
	{
		sosRoute.POST("/emergency", controller.TriggerEmergency)
		sosRoute.POST("/:sosId/accept", controller.AcceptSOS)
		sosRoute.PUT("/:sosId/resolve", controller.ResolveSOS)
		sosRoute.GET("/current", controller.GetCurrentSOS)
	}
}
