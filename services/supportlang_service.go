package services

import (
	"errors"
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/repositories"
)

type SupportLangService struct {
	repo *repositories.SupportLangRepository
}

func NewSupportLangService(repo *repositories.SupportLangRepository) *SupportLangService {
	return &SupportLangService{
		repo: repo,
	}
}

func (s *SupportLangService) GetAllSupportLangs() ([]models.SupportLang, error) {
	return s.repo.GetAll()
}

func (s *SupportLangService) GetSupportLangByID(id uint) (models.SupportLang, error) {
	if id == 0 {
		return models.SupportLang{}, errors.New("無效的語言ID")
	}
	return s.repo.GetByID(id)
}

func (s *SupportLangService) GetSupportLangByName(name string) (models.SupportLang, error) {
	if name == "" {
		return models.SupportLang{}, errors.New("語言名稱不能為空")
	}
	return s.repo.GetByName(name)
}

func (s *SupportLangService) GetActiveLanguages() ([]models.SupportLang, error) {
	return s.repo.GetActiveLanguages()
}

func (s *SupportLangService) CreateSupportLang(lang models.SupportLang) (models.SupportLang, error) {
	if lang.Name == "" {
		return models.SupportLang{}, errors.New("語言名稱不能為空")
	}

	if lang.RunCommand == "" {
		return models.SupportLang{}, errors.New("執行命令不能為空")
	}

	if lang.FileExtension == "" {
		return models.SupportLang{}, errors.New("檔案副檔名不能為空")
	}

	existingLang, err := s.repo.GetByName(lang.Name)
	if err == nil && existingLang.ID != 0 {
		return models.SupportLang{}, errors.New("語言已存在")
	}

	return s.repo.Create(lang)
}

func (s *SupportLangService) UpdateSupportLang(id uint, langUpdate models.SupportLang) error {
	if id == 0 {
		return errors.New("無效的語言ID")
	}

	existingLang, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("語言不存在")
	}

	if langUpdate.Name != "" && langUpdate.Name != existingLang.Name {
		duplicateLang, err := s.repo.GetByName(langUpdate.Name)
		if err == nil && duplicateLang.ID != 0 && duplicateLang.ID != id {
			return errors.New("語言名稱已存在")
		}
		existingLang.Name = langUpdate.Name
	}

	if langUpdate.Version != "" {
		existingLang.Version = langUpdate.Version
	}

	if langUpdate.CompileCommand != "" {
		existingLang.CompileCommand = langUpdate.CompileCommand
	}

	if langUpdate.RunCommand != "" {
		existingLang.RunCommand = langUpdate.RunCommand
	}

	if langUpdate.FileExtension != "" {
		existingLang.FileExtension = langUpdate.FileExtension
	}

	existingLang.IsActive = langUpdate.IsActive

	return s.repo.Update(existingLang)
}

func (s *SupportLangService) DeleteSupportLang(id uint) error {
	if id == 0 {
		return errors.New("無效的語言ID")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("語言不存在")
	}

	return s.repo.Delete(id)
}

func (s *SupportLangService) ToggleLanguageActive(id uint, isActive bool) error {
	if id == 0 {
		return errors.New("無效的語言ID")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("語言不存在")
	}

	return s.repo.ToggleActive(id, isActive)
}
