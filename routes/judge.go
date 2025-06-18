package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupJudgeRoutes(router *gin.RouterGroup, judgeController *controllers.JudgeController, solutionController *controllers.SolutionController) {

	publicJudges := router.Group("/judges")
	{
		publicJudges.GET("", judgeController.GetJudges)
		publicJudges.GET("/:id", judgeController.GetJudge)
		publicJudges.GET("/:id/code/:language", judgeController.GetJudgeCode)

		// 根據標籤獲取題目
		publicJudges.GET("/tags", judgeController.GetJudgesByTags)

		// 根據難度獲取題目
		publicJudges.GET("/difficulty/:difficulty", judgeController.GetJudgesByDifficulty)

		// 獲取特定題目的所有答案
		publicJudges.GET("/:id/solutions", solutionController.GetJudgeSolutions)
	}

	protectedJudges := router.Group("/judges", middleware.JwtAuth())
	{
		protectedJudges.POST("", judgeController.CreateJudge)
		protectedJudges.PUT("/:id", judgeController.UpdateJudge)
		protectedJudges.DELETE("/:id", judgeController.DeleteJudge)
	}

}
