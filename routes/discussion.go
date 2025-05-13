package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupDiscussionRoutes(router *gin.RouterGroup, controller *controllers.DiscussionController) {
	publicDiscussions := router.Group("/discussions")
	{
		publicDiscussions.GET("", controller.GetDiscussions)
		publicDiscussions.GET("/:id", controller.GetDiscussion)

		// 獲取特定用戶的討論
		publicDiscussions.GET("/user/:user_id", controller.GetUserDiscussions)

		// 獲取子討論（回覆）
		publicDiscussions.GET("/parent/:parent_type/:parent_id", controller.GetChildDiscussions)

		// 搜索討論
		publicDiscussions.GET("/search", controller.SearchDiscussions)
	}

	protectedDiscussions := router.Group("/discussions", middleware.JwtAuth())
	{
		protectedDiscussions.POST("", controller.CreateDiscussion)
		protectedDiscussions.PUT("/:id", controller.UpdateDiscussion)
		protectedDiscussions.DELETE("/:id", controller.DeleteDiscussion)
	}
}
