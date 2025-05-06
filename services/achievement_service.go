package services

import (
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/repositories"
)

type AchievementService struct {
	repo *repositories.AchievementRepository
}

func NewAchievementService(repo *repositories.AchievementRepository) *AchievementService {
	return &AchievementService{
		repo: repo,
	}
}

func (s *AchievementService) GetAllAchievements() ([]models.Achievement, error) {
	return s.repo.GetAll()
}

func (s *AchievementService) GetAchievementByID(id uint) (models.Achievement, error) {
	return s.repo.GetByID(id)
}

func (s *AchievementService) CreateAchievement(achievement models.Achievement) (models.Achievement, error) {
	return s.repo.Create(achievement)
}

func (s *AchievementService) UpdateAchievement(achievement models.Achievement) error {
	return s.repo.Update(achievement)
}

func (s *AchievementService) DeleteAchievement(id uint) error {
	return s.repo.Delete(id)
}
