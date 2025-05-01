package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupSolutionRoutes(router *gin.RouterGroup, solutionController *controllers.SolutionController) {
	solutions := router.Group("/solutions", middleware.JwtAuth())
	{
		solutions.GET("", solutionController.GetSolutions)
		solutions.GET("/:id", solutionController.GetSolution)
		solutions.POST("", solutionController.SubmitSolution)
		solutions.PUT("/:id/result", solutionController.UpdateSolutionResult)
		solutions.DELETE("/:id", solutionController.DeleteSolution)

		solutions.GET("/best", solutionController.GetBestSolutions)
	}
}
