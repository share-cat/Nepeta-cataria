package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupSocialRoutes(router *gin.RouterGroup, controller *controllers.SocialController) {

	publicSocials := router.Group("/socials")
	{
		publicSocials.GET("", controller.GetSocials)
		publicSocials.GET("/:id", controller.GetSocial)
	}

	protectedSocials := router.Group("/socials", middleware.JwtAuth())
	{
		protectedSocials.POST("", controller.CreateSocial)
		protectedSocials.PUT("/:id", controller.UpdateSocial)
		protectedSocials.DELETE("/:id", controller.DeleteSocial)
	}
}
