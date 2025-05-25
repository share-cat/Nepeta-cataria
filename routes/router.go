package routes

import (
	"share-cat/Nepeta-cataria/boot"
	"share-cat/Nepeta-cataria/docs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(g *gin.Engine, db *gorm.DB, controllers boot.Controllers) {

	api := g.Group("/api")
	docs.SwaggerInfo.BasePath = "/api"

	SetupAchievementRoutes(api, controllers.AchievementController)
	SetupChatSessionRoutes(api, controllers.ChatSessionController)
	SetupDiscussionRoutes(api, controllers.DiscussionController)
	SetupJudgeRoutes(api, controllers.JudgeController, controllers.SolutionController)
	SetupSchoolRoutes(api, controllers.SchoolController)
	SetupSocialRoutes(api, controllers.SocialController)
	SetupSolutionRoutes(api, controllers.SolutionController)
	SetupSubmitEffectRoutes(api, controllers.SubmitEffectController)
	SetupUserProfileRoutes(api, controllers.UserProfileController)
	SetupUserRoutes(api, controllers.UserController, controllers.JudgeController, controllers.SolutionController)
	SetupAuthRoutes(api, controllers.AuthController)
	SetupLoadjudgeRoutes(api, controllers.LoadJudgeController)
	SetupSwaggerRoutes(api)

}
