// routes/chat_session.go
package routes

import (
	"share-cat/Nepeta-cataria/controllers"

	"github.com/gin-gonic/gin"
)

func SetupChatSessionRoutes(router *gin.RouterGroup, controller *controllers.ChatSessionController) {
	chatSessions := router.Group("/chat-sessions")
	{
		chatSessions.GET("", controller.GetChatSessions)
		chatSessions.GET("/:id", controller.GetChatSession)
		chatSessions.POST("", controller.CreateChatSession)
		chatSessions.PUT("/:id", controller.UpdateChatSession)
		chatSessions.DELETE("/:id", controller.DeleteChatSession)

		// 更新聊天歷史
		chatSessions.PUT("/:id/history", controller.UpdateChatHistory)

		// 更新標題
		chatSessions.PUT("/:id/title", controller.UpdateTitle)

		// 根據用戶ID獲取聊天會話
		chatSessions.GET("/user/:user_id", controller.GetUserChatSessions)

		// 根據評審ID獲取聊天會話
		chatSessions.GET("/judge/:judge_id", controller.GetJudgeChatSessions)

		// 獲取特定用戶和特定評審的聊天會話
		chatSessions.GET("/user/:user_id/judge/:judge_id", controller.GetUserJudgeChatSessions)
	}
}
