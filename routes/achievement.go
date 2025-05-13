package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAchievementRoutes(router *gin.RouterGroup, controller *controllers.AchievementController) {

	publicAchievements := router.Group("/achievements")
	{
		publicAchievements.GET("", controller.GetAchievements)
		publicAchievements.GET("/:id", controller.GetAchievement)
	}

	protectedAchievements := router.Group("/achievements", middleware.JwtAuth())
	{
		protectedAchievements.POST("", controller.CreateAchievement)
		protectedAchievements.PUT("", controller.UpdateAchievement)
		protectedAchievements.DELETE("/:id", controller.DeleteAchievement)
	}
}
