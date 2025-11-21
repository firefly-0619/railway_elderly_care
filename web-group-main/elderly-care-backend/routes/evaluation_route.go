package routes

import (
	"elderly-care-backend/controllers"
	"github.com/gin-gonic/gin"
)

func EvaluationRoute(e *gin.Engine) {
	evaluationController := controllers.EvaluationController{}
	evaluationRoute := e.Group("/evaluation")
	{
		evaluationRoute.PUT("/account", evaluationController.EvaluationAccount)
		evaluationRoute.GET("/account", evaluationController.GetAccountEvaluation)
	}

}
