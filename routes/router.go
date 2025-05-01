// routes/router.go
package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/repositories"
	"share-cat/Nepeta-cataria/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(g *gin.Engine, db *gorm.DB) {

	userRepo := repositories.NewUserRepository(db)
	judgeRepo := repositories.NewJudgeRepository(db)
	solutionRepo := repositories.NewSolutionRepository(db)
	langRepo := repositories.NewSupportLangRepository(db)

	userService := services.NewUserService(userRepo)
	judgeService := services.NewJudgeService(judgeRepo, userRepo)
	solutionService := services.NewSolutionService(solutionRepo, userRepo, judgeRepo, langRepo)

	userController := controllers.NewUserController(userService)
	judgeController := controllers.NewJudgeController(judgeService)
	solutionController := controllers.NewSolutionController(solutionService)

	api := g.Group("/api")

	SetupUserRoutes(api, userController, judgeController, solutionController)
	SetupJudgeRoutes(api, judgeController, solutionController)
	SetupSolutionRoutes(api, solutionController)

}
