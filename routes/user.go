package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, userController *controllers.UserController, judgeController *controllers.JudgeController, solutionController *controllers.SolutionController) {

	publicUsers := router.Group("/users")
	{
		publicUsers.GET("", userController.GetUsers)
		publicUsers.GET("/:id", userController.GetUser)

		// 獲取特定用戶創建的題目
		publicUsers.GET("/:id/judges", judgeController.GetJudgesByCreator)

		// 獲取特定用戶的答案
		publicUsers.GET("/:id/solutions", solutionController.GetUserSolutions)
	}

	protectedUsers := router.Group("/users", middleware.JwtAuth())
	{
		protectedUsers.POST("", userController.CreateUser)
		protectedUsers.PUT("/:id", userController.UpdateUser)
		protectedUsers.PUT("/:id/password", userController.UpdatePassword)
		protectedUsers.DELETE("/:id", userController.DeleteUser)
		protectedUsers.GET("/:id/profile", userController.GetUserProfile)
	}
}
