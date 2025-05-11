package routes

import (
	"share-cat/Nepeta-cataria/controllers"

	"github.com/gin-gonic/gin"
)

func SetupUserProfileRoutes(router *gin.RouterGroup, controller *controllers.UserProfileController) {
	userProfiles := router.Group("/user_profiles")
	{
		userProfiles.GET("", controller.GetAllUserProfiles)
		userProfiles.GET("/:id", controller.GetUserProfileByID)
		userProfiles.GET("/user/:id", controller.GetUserProfileByUserID)
		userProfiles.POST("", controller.CreateUserProfile)
		userProfiles.PUT("/:id", controller.UpdateUserProfile)
		userProfiles.PUT("/:user_id/avatar", controller.UpdateAvatarURL)
		userProfiles.DELETE("/:id", controller.DeleteUserProfile)
	}
}
