package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserProfileRoutes(router *gin.RouterGroup, controller *controllers.UserProfileController) {
	userProfiles := router.Group("/user_profiles")
	{
		userProfiles.GET("", controller.GetAllUserProfiles)
		userProfiles.GET("/:id", controller.GetUserProfileByID)
		userProfiles.GET("/user/:id", controller.GetUserProfileByUserID)
	}

	protectedProfiles := router.Group("/user_profiles", middleware.JwtAuth())
	{
		protectedProfiles.POST("", controller.CreateUserProfile)
		protectedProfiles.PUT("/:id", controller.UpdateUserProfile)
		protectedProfiles.PUT("/:id/avatar", controller.UpdateAvatarURL)
		protectedProfiles.DELETE("/:id", controller.DeleteUserProfile)
	}
}
