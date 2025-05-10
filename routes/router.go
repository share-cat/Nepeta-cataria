// routes/router.go
package routes

import (
	"share-cat/Nepeta-cataria/boot"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(g *gin.Engine, db *gorm.DB, controllers boot.Controllers) {

	api := g.Group("/api")

	SetupUserRoutes(api, controllers.UserController, controllers.JudgeController, controllers.SolutionController)
	SetupJudgeRoutes(api, controllers.JudgeController, controllers.SolutionController)
	SetupSolutionRoutes(api, controllers.SolutionController)

}
