package boot

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/repositories"
	"share-cat/Nepeta-cataria/services"

	"gorm.io/gorm"
)

type Repositories struct {
	AchievementRepo  *repositories.AchievementRepository
	ChatSessionRepo  *repositories.ChatSessionRepository
	DiscussionRepo   *repositories.DiscussionRepository
	JudgeRepo        *repositories.JudgeRepository
	SchoolRepo       *repositories.SchoolRepository
	SocialRepo       *repositories.SocialRepository
	SolutionRepo     *repositories.SolutionRepository
	SolutionTaskRepo *repositories.SolutionTaskRepository
	SubmiteffectRepo *repositories.SubmitEffectRepository
	SupportLangRepo  *repositories.SupportLangRepository
	UserRepo         *repositories.UserRepository
	UserProfileRepo  *repositories.UserProfileRepository
}

type Services struct {
	AchievementService  *services.AchievementService
	ChatSessionService  *services.ChatSessionService
	DiscussionService   *services.DiscussionService
	LoadJudgeService    *services.LoadJudgeService
	JudgeService        *services.JudgeService
	SchoolService       *services.SchoolService
	SocialService       *services.SocialService
	SolutionService     *services.SolutionService
	SolutionTaskService *services.SolutionTaskService
	SubmitEffectService *services.SubmitEffectService
	SupportLangService  *services.SupportLangService
	UserService         *services.UserService
	UserProfileService  *services.UserProfileService
}

type Controllers struct {
	AuthController         *controllers.AuthController
	AchievementController  *controllers.AchievementController
	ChatSessionController  *controllers.ChatSessionController
	DiscussionController   *controllers.DiscussionController
	LoadJudgeController    *controllers.LoadJudgeController
	SchoolController       *controllers.SchoolController
	SocialController       *controllers.SocialController
	SubmitEffectController *controllers.SubmitEffectController
	UserProfileController  *controllers.UserProfileController
	UserController         *controllers.UserController
	JudgeController        *controllers.JudgeController
	SolutionController     *controllers.SolutionController
	SolutionTaskController *controllers.SolutionTaskController
}

type Application struct {
	Repositories Repositories
	Services     Services
	Controllers  Controllers
}

func NewRepositories(db *gorm.DB) Repositories {
	return Repositories{
		AchievementRepo:  repositories.NewAchievementRepository(db),
		ChatSessionRepo:  repositories.NewChatSessionRepository(db),
		DiscussionRepo:   repositories.NewDiscussionRepository(db),
		JudgeRepo:        repositories.NewJudgeRepository(db),
		SchoolRepo:       repositories.NewSchoolRepository(db),
		SocialRepo:       repositories.NewSocialRepository(db),
		SolutionRepo:     repositories.NewSolutionRepository(db),
		SolutionTaskRepo: repositories.NewSolutionTaskRepository(db),
		SubmiteffectRepo: repositories.NewSubmitEffectRepository(db),
		SupportLangRepo:  repositories.NewSupportLangRepository(db),
		UserRepo:         repositories.NewUserRepository(db),
		UserProfileRepo:  repositories.NewUserProfileRepository(db),
	}
}

func NewServices(repos Repositories) Services {
	return Services{
		AchievementService:  services.NewAchievementService(repos.AchievementRepo),
		ChatSessionService:  services.NewChatSessionService(repos.ChatSessionRepo, repos.UserRepo, repos.JudgeRepo),
		DiscussionService:   services.NewDiscussionService(repos.DiscussionRepo, repos.UserRepo),
		LoadJudgeService:    services.NewLoadJudgeService(repos.JudgeRepo),
		JudgeService:        services.NewJudgeService(repos.JudgeRepo, repos.UserRepo),
		SchoolService:       services.NewSchoolService(repos.SchoolRepo),
		SocialService:       services.NewSocialService(repos.SocialRepo),
		SolutionService:     services.NewSolutionService(repos.SolutionRepo, repos.UserRepo, repos.JudgeRepo, repos.SupportLangRepo),
		SolutionTaskService: services.NewSolutionTaskService(repos.SolutionTaskRepo, repos.SolutionRepo),
		SubmitEffectService: services.NewSubmitEffectService(repos.SubmiteffectRepo),
		SupportLangService:  services.NewSupportLangService(repos.SupportLangRepo),
		UserService:         services.NewUserService(repos.UserRepo),
		UserProfileService:  services.NewUserProfileService(repos.UserProfileRepo, repos.UserRepo),
	}
}

func NewControllers(services Services) Controllers {
	return Controllers{
		AuthController:         controllers.NewAuthController(services.UserService),
		AchievementController:  controllers.NewAchievementController(services.AchievementService),
		ChatSessionController:  controllers.NewChatSessionController(services.ChatSessionService),
		DiscussionController:   controllers.NewDiscussionController(services.DiscussionService),
		LoadJudgeController:    controllers.NewLoadJudgeController(services.LoadJudgeService),
		SchoolController:       controllers.NewSchoolController(services.SchoolService),
		SocialController:       controllers.NewSocialController(services.SocialService),
		SubmitEffectController: controllers.NewSubmitEffectController(services.SubmitEffectService),
		UserProfileController:  controllers.NewUserProfileController(services.UserProfileService),
		UserController:         controllers.NewUserController(services.UserService),
		JudgeController:        controllers.NewJudgeController(services.JudgeService),
		SolutionController:     controllers.NewSolutionController(services.SolutionService, services.SolutionTaskService, services.JudgeService, services.SupportLangService),
		SolutionTaskController: controllers.NewSolutionTaskController(services.SolutionTaskService),
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
