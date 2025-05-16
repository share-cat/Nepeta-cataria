package routes

import (
	"share-cat/Nepeta-cataria/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.RouterGroup, controller *controllers.AuthController) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/refresh-token", controller.RefreshToken)
	}
}
