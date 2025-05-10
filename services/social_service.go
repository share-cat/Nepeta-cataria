package services

import (
	"errors"
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/repositories"
)

type SocialService struct {
	repo *repositories.SociaRepository
}

func NewSocialService(repo *repositories.SociaRepository) *SocialService {
	return &SocialService{
		repo: repo,
	}
}

func (s *SocialService) GetAllSocials() ([]models.Social, error) {
	return s.repo.GetAll()
}

func (s *SocialService) GetSocialByID(id uint) (models.Social, error) {
	if id == 0 {
		return models.Social{}, errors.New("invalid social ID")
	}
	return s.repo.GetByID(id)
}

func (s *SocialService) CreateSocial(social models.Social) (models.Social, error) {
	return s.repo.Create(social)
}

func (s *SocialService) UpdateSocial(social models.Social) error {
	if social.ID == 0 {
		return errors.New("invalid social ID")
	}
	return s.repo.Update(social)
}

func (s *SocialService) DeleteSocial(id uint) error {
	if id == 0 {
		return errors.New("invalid social ID")
	}
	return s.repo.Delete(id)
}
