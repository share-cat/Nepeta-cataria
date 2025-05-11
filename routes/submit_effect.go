package routes

import (
	"share-cat/Nepeta-cataria/controllers"

	"github.com/gin-gonic/gin"
)

func SetupSubmitEffectRoutes(router *gin.RouterGroup, controller *controllers.SubmitEffectController) {
	submitEffects := router.Group("/submit_effects")
	{
		submitEffects.GET("", controller.GetSubmitEffects)
		submitEffects.GET("/:id", controller.GetSubmitEffect)
		submitEffects.POST("", controller.CreateSubmitEffect)
		submitEffects.PUT("/:id", controller.UpdateSubmitEffect)
		submitEffects.DELETE("/:id", controller.DeleteSubmitEffect)
	}
}
