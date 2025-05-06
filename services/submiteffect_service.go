package services

import (
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/repositories"
)

type SubmitEffectService struct {
	repo *repositories.SubmitEffectRepository
}

func NewSubmitEffectService(repo *repositories.SubmitEffectRepository) *SubmitEffectService {
	return &SubmitEffectService{
		repo: repo,
	}
}

func (s *SubmitEffectService) GetAllSubmitEffects() ([]models.SubmitEffect, error) {
	return s.repo.GetAll()
}

func (s *SubmitEffectService) GetSubmitEffectByID(id uint) (models.SubmitEffect, error) {
	return s.repo.GetByID(id)
}

func (s *SubmitEffectService) SearchSubmitEffectByName(name string) (models.SubmitEffect, error) {
	return s.repo.SearchByName(name)
}

func (s *SubmitEffectService) CreateSubmitEffect(submitEffect models.SubmitEffect) (models.SubmitEffect, error) {
	return s.repo.Create(submitEffect)
}

func (s *SubmitEffectService) UpdateSubmitEffect(submitEffect models.SubmitEffect) error {
	return s.repo.Update(submitEffect)
}

func (s *SubmitEffectService) DeleteSubmitEffect(id uint) error {
	return s.repo.Delete(id)
}
