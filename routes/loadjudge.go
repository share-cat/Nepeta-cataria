package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupLoadjudgeRoutes(router *gin.RouterGroup, controller *controllers.LoadJudgeController) {

	protectedSocials := router.Group("/load-judge", middleware.JwtAuth())
	{
		protectedSocials.POST("/upload", controller.UploadAndLoadJudge)
		protectedSocials.POST("/path", controller.LoadJudgeFromPath)
		protectedSocials.POST("/validate", controller.ValidateJudgeZip)
		protectedSocials.GET("/directory/:id", controller.GetJudgeDirectory)
		protectedSocials.GET("/download/:id", controller.DownloadJudgeFile)
		protectedSocials.GET("/base-directory", controller.GetJudgesBaseDirectory)
	}
}
