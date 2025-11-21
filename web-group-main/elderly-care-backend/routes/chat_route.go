package routes

import (
	"elderly-care-backend/controllers"
	"github.com/gin-gonic/gin"
)

func ChatRoute(e *gin.Engine) {
	chatManager := controllers.NewConnectionManager()
	go chatManager.Start()
	chatRoute := e.Group("/chat")
	{
		chatRoute.GET("", chatManager.HandleWebSocket)
		chatRoute.GET("/record", chatManager.GetChatRecord)
		chatRoute.GET("/contactList", chatManager.GetRecentlyChatList)
	}
}
