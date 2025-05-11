package routes

import (
	"share-cat/Nepeta-cataria/controllers"

	"github.com/gin-gonic/gin"
)

func SetupDiscussionRoutes(router *gin.RouterGroup, controller *controllers.DiscussionController) {
	discussions := router.Group("/discussions")
	{
		discussions.GET("", controller.GetDiscussions)
		discussions.GET("/:id", controller.GetDiscussion)
		discussions.POST("", controller.CreateDiscussion)
		discussions.PUT("/:id", controller.UpdateDiscussion)
		discussions.DELETE("/:id", controller.DeleteDiscussion)

		// 獲取特定用戶的討論
		discussions.GET("/user/:user_id", controller.GetUserDiscussions)

		// 獲取子討論（回覆）
		discussions.GET("/parent/:parent_type/:parent_id", controller.GetChildDiscussions)

		// 搜索討論
		discussions.GET("/search", controller.SearchDiscussions)
	}
}
