package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, userController *controllers.UserController, judgeController *controllers.JudgeController, solutionController *controllers.SolutionController) {
	users := router.Group("/users", middleware.JwtAuth())
	{
		users.GET("", userController.GetUsers)
		users.GET("/:id", userController.GetUser)
		users.POST("", userController.CreateUser)
		users.PUT("/:id", userController.UpdateUser)
		users.PUT("/:id/password", userController.UpdatePassword)
		users.DELETE("/:id", userController.DeleteUser)
		users.GET("/:id/profile", userController.GetUserProfile)

		// 獲取特定用戶創建的題目
		users.GET("/:id/judges", judgeController.GetJudgesByCreator)

		// 獲取特定用戶的答案
		users.GET("/:id/solutions", solutionController.GetUserSolutions)
	}
}
