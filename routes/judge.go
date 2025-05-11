package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupJudgeRoutes(router *gin.RouterGroup, judgeController *controllers.JudgeController, solutionController *controllers.SolutionController) {
	judges := router.Group("/judges", middleware.JwtAuth())
	{
		judges.GET("", judgeController.GetJudges)
		judges.GET("/:id", judgeController.GetJudge)
		judges.POST("", judgeController.CreateJudge)
		judges.PUT("/:id", judgeController.UpdateJudge)
		judges.DELETE("/:id", judgeController.DeleteJudge)

		// 根據標籤獲取題目
		judges.GET("/tags", judgeController.GetJudgesByTags)

		// 根據難度獲取題目
		judges.GET("/difficulty/:difficulty", judgeController.GetJudgesByDifficulty)

		// 獲取特定題目的所有答案
		judges.GET("/:judge_id/solutions", solutionController.GetJudgeSolutions)
	}

}
