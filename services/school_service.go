package services

import (
	"errors"
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/repositories"
)

type SchoolService struct {
	repo *repositories.SchoolRepository
}

func NewSchoolService(repo *repositories.SchoolRepository) *SchoolService {
	return &SchoolService{
		repo: repo,
	}
}

func (s *SchoolService) GetAllSchools() ([]models.School, error) {
	return s.repo.GetAll()
}

func (s *SchoolService) GetSchoolByID(id uint) (models.School, error) {
	if id == 0 {
		return models.School{}, errors.New("invalid school ID")
	}
	return s.repo.GetByID(id)
}

func (s *SchoolService) GetSchoolByName(name string) (models.School, error) {
	return s.repo.SearchByName(name)
}

func (s *SchoolService) CreateSchool(school models.School) (models.School, error) {
	result, err := s.repo.SearchByName(school.SchoolName)
	if err == nil {
		return models.School{}, errors.New("學校已存在")
	}
	if result.ID != 0 {
		return models.School{}, errors.New("學校已存在")
	}
	return s.repo.Create(school)

}

func (s *SchoolService) UpdateSchool(school models.School) error {
	if school.ID == 0 {
		return errors.New("invalid school ID")
	}
	return s.repo.Update(school)
}

func (s *SchoolService) DeleteSchool(id uint) error {
	if id == 0 {
		return errors.New("invalid school ID")
	}
	return s.repo.Delete(id)
}
