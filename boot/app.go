package boot

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/repositories"
	"share-cat/Nepeta-cataria/services"

	"gorm.io/gorm"
)

type Repositories struct {
	AchievementRepo  *repositories.AchievementRepository
	ChatsessionRepo  *repositories.ChatSessionRepository
	DiscussionRepo   *repositories.DiscussionRepository
	JudgeRepo        *repositories.JudgeRepository
	SchoolRepo       *repositories.SchoolRepository
	SocialRepo       *repositories.SocialRepository
	SolutionRepo     *repositories.SolutionRepository
	SubmiteffectRepo *repositories.SubmitEffectRepository
	SupportLangRepo  *repositories.SupportLangRepository
	UserRepo         *repositories.UserRepository
	UserProfileRepo  *repositories.UserProfileRepository
}

type Services struct {
	AchievementService  *services.AchievementService
	ChatSessionService  *services.ChatSessionService
	DiscussionService   *services.DiscussionService
	JudgeService        *services.JudgeService
	SchoolService       *services.SchoolService
	SocialService       *services.SocialService
	SolutionService     *services.SolutionService
	SubmitEffectService *services.SubmitEffectService
	UserService         *services.UserService
	UserProfileService  *services.UserProfileService
}

type Controllers struct {
	AchievementController  *controllers.AchievementController
	ChatSessionController  *controllers.ChatSessionController
	DiscussionController   *controllers.DiscussionController
	SchoolController       *controllers.SchoolController
	SocialController       *controllers.SocialController
	SubmitEffectController *controllers.SubmitEffectController
	UserProfileController  *controllers.UserProfileController
	UserController         *controllers.UserController
	JudgeController        *controllers.JudgeController
	SolutionController     *controllers.SolutionController
}

type Application struct {
	Repositories Repositories
	Services     Services
	Controllers  Controllers
}

func NewRepositories(db *gorm.DB) Repositories {
	return Repositories{
		UserRepo:        repositories.NewUserRepository(db),
		JudgeRepo:       repositories.NewJudgeRepository(db),
		SolutionRepo:    repositories.NewSolutionRepository(db),
		SupportLangRepo: repositories.NewSupportLangRepository(db),
	}
}

func NewServices(repos Repositories) Services {
	return Services{
		UserService:     services.NewUserService(repos.UserRepo),
		JudgeService:    services.NewJudgeService(repos.JudgeRepo, repos.UserRepo),
		SolutionService: services.NewSolutionService(repos.SolutionRepo, repos.UserRepo, repos.JudgeRepo, repos.SupportLangRepo),
	}
}

func NewControllers(services Services) Controllers {
	return Controllers{
		UserController:     controllers.NewUserController(services.UserService),
		JudgeController:    controllers.NewJudgeController(services.JudgeService),
		SolutionController: controllers.NewSolutionController(services.SolutionService),
	}
}

func NewApplication(db *gorm.DB) Application {
	repositories := NewRepositories(db)
	services := NewServices(repositories)
	controllers := NewControllers(services)

	return Application{
		Repositories: repositories,
		Services:     services,
		Controllers:  controllers,
	}
}
