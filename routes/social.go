package routes

import (
	"share-cat/Nepeta-cataria/controllers"

	"github.com/gin-gonic/gin"
)

func SetupSocialRoutes(router *gin.RouterGroup, controller *controllers.SocialController) {
	socials := router.Group("/socials")
	{
		socials.GET("", controller.GetSocials)
		socials.GET("/:id", controller.GetSocial)
		socials.POST("", controller.CreateSocial)
		socials.PUT("/:id", controller.UpdateSocial)
		socials.DELETE("/:id", controller.DeleteSocial)
	}
}
