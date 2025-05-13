// routes/chat_session.go
package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupChatSessionRoutes(router *gin.RouterGroup, controller *controllers.ChatSessionController) {

	publicChatSessions := router.Group("/chat-sessions")
	{
		publicChatSessions.GET("", controller.GetChatSessions)
		publicChatSessions.GET("/:id", controller.GetChatSession)

		// 根據用戶ID獲取聊天會話
		publicChatSessions.GET("/user/:user_id", controller.GetUserChatSessions)

		// 根據評審ID獲取聊天會話
		publicChatSessions.GET("/judge/:judge_id", controller.GetJudgeChatSessions)

		// 獲取特定用戶和特定評審的聊天會話
		publicChatSessions.GET("/user/:user_id/judge/:judge_id", controller.GetUserJudgeChatSessions)
	}

	protectedChatSessions := router.Group("/chat-sessions", middleware.JwtAuth())
	{
		protectedChatSessions.POST("", controller.CreateChatSession)
		protectedChatSessions.PUT("/:id", controller.UpdateChatSession)
		protectedChatSessions.DELETE("/:id", controller.DeleteChatSession)

		// 更新聊天歷史
		protectedChatSessions.PUT("/:id/history", controller.UpdateChatHistory)

		// 更新標題
		protectedChatSessions.PUT("/:id/title", controller.UpdateTitle)
	}
}
