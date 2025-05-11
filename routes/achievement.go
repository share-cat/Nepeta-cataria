package routes

import (
	"share-cat/Nepeta-cataria/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAchievementRoutes(router *gin.RouterGroup, controller *controllers.AchievementController) {
	achievements := router.Group("/achievements")
	{
		achievements.GET("", controller.GetAchievements)
		achievements.GET("/:id", controller.GetAchievement)
		achievements.POST("", controller.CreateAchievement)
		achievements.PUT("", controller.UpdateAchievement)
		achievements.DELETE("/:id", controller.DeleteAchievement)
	}
}
