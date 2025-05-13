package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupSubmitEffectRoutes(router *gin.RouterGroup, controller *controllers.SubmitEffectController) {

	publicSubmitEffects := router.Group("/submit_effects")
	{
		publicSubmitEffects.GET("", controller.GetSubmitEffects)
		publicSubmitEffects.GET("/:id", controller.GetSubmitEffect)
	}

	protectedSubmitEffects := router.Group("/submit_effects", middleware.JwtAuth())
	{
		protectedSubmitEffects.POST("", controller.CreateSubmitEffect)
		protectedSubmitEffects.PUT("/:id", controller.UpdateSubmitEffect)
		protectedSubmitEffects.DELETE("/:id", controller.DeleteSubmitEffect)
	}
}
