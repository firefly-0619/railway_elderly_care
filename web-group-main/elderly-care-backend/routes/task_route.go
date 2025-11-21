package routes

import (
	"elderly-care-backend/controllers"

	"github.com/gin-gonic/gin"
)

func TaskRoute(e *gin.Engine) {
	controller := controllers.NewTaskController()
	taskRoute := e.Group("/tasks")
	{
		taskRoute.POST("/create", controller.CreateTask)
		taskRoute.GET("/nearby", controller.GetNearbyTasks)
		taskRoute.POST("/:taskId/accept", controller.AcceptTask)
	}
}
