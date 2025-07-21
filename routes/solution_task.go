package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupSolutionTaskRoutes(router *gin.Engine, controller *controllers.SolutionTaskController) {
	solutionTaskRoutes := router.Group("/api/v1")
	{
		// 獲取解答的所有任務
		solutionTaskRoutes.GET("/solutions/:solution_id/tasks", controller.GetSolutionTasks)

		// 獲取解答匯總資訊
		solutionTaskRoutes.GET("/solutions/:solution_id/summary", controller.GetSolutionSummary)

		// 獲取特定任務
		solutionTaskRoutes.GET("/solution-tasks/:id", controller.GetSolutionTask)

		// 需要認證的路由
		authRoutes := solutionTaskRoutes.Group("")
		authRoutes.Use(middleware.JwtAuth())
		{
			// 更新任務狀態
			authRoutes.PUT("/solution-tasks/:id/status", controller.UpdateTaskStatus)

			// 刪除任務
			authRoutes.DELETE("/solution-tasks/:id", controller.DeleteTask)
		}
	}
}
