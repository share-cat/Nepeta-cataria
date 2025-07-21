package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupSolutionRoutes(router *gin.RouterGroup, solutionController *controllers.SolutionController) {
	publicSolutions := router.Group("/solutions")
	{
		publicSolutions.GET("", solutionController.GetSolutions)
		publicSolutions.GET("/:id", solutionController.GetSolution)
		publicSolutions.GET("/best", solutionController.GetBestSolutions)
	}

	protectedSolutions := router.Group("/solutions", middleware.JwtAuth())
	{
		protectedSolutions.POST("", solutionController.SubmitSolution)
		protectedSolutions.DELETE("/:id", solutionController.DeleteSolution)
	}
}
