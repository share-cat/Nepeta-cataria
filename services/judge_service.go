package services

import (
	"errors"

	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/repositories"
)

type JudgeService struct {
	repo     *repositories.JudgeRepository
	userRepo *repositories.UserRepository
}

func NewJudgeService(repo *repositories.JudgeRepository, userRepo *repositories.UserRepository) *JudgeService {
	return &JudgeService{
		repo:     repo,
		userRepo: userRepo,
	}
}

// GetAllJudges 獲取所有題目
func (s *JudgeService) GetAllJudges() ([]models.Judge, error) {
	return s.repo.GetAll()
}

// GetJudgesWithFilters 使用篩選條件獲取題目
func (s *JudgeService) GetJudgesWithFilters(search, judgeType string, page int, activeOnly bool) ([]models.Judge, int, error) {
	return s.repo.GetJudgesWithFilters(search, judgeType, page, activeOnly)
}

// GetJudgeByID 根據ID獲取題目
func (s *JudgeService) GetJudgeByID(id uint) (models.Judge, error) {
	return s.repo.GetByID(id)
}

// GetActiveJudges 獲取所有活躍的題目
func (s *JudgeService) GetActiveJudges() ([]models.Judge, error) {
	return s.repo.GetActiveJudges()
}

// CreateJudge 創建新題目
func (s *JudgeService) CreateJudge(judge models.Judge) (models.Judge, error) {

	_, err := s.userRepo.GetByID(uint(judge.CreatorID))
	if err != nil {
		return models.Judge{}, errors.New("創建者不存在")
	}

	if judge.Type == "" {
		judge.Type = "一般題目"
	}

	judge.IsActive = true

	return s.repo.Create(judge)
}

// UpdateJudge 更新題目
func (s *JudgeService) UpdateJudge(id uint, judgeUpdate models.Judge) error {
	judge, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if judgeUpdate.Name != "" {
		judge.Name = judgeUpdate.Name
	}

	if judgeUpdate.Description != "" {
		judge.Description = judgeUpdate.Description
	}

	if len(judgeUpdate.Tags) > 0 {
		judge.Tags = judgeUpdate.Tags
	}

	if judgeUpdate.Type != "" {
		judge.Type = judgeUpdate.Type
	}

	if judgeUpdate.Difficulty > 0 {
		judge.Difficulty = judgeUpdate.Difficulty
	}

	judge.IsActive = judgeUpdate.IsActive

	return s.repo.Update(judge)
}

// DeleteJudge 刪除題目
func (s *JudgeService) DeleteJudge(id uint) error {
	return s.repo.Delete(id)
}

// GetJudgesByCreator 獲取特定用戶創建的題目
func (s *JudgeService) GetJudgesByCreator(creatorID uint) ([]models.Judge, error) {

	_, err := s.userRepo.GetByID(creatorID)
	if err != nil {
		return nil, errors.New("用戶不存在")
	}

	return s.repo.GetJudgesByCreator(creatorID)
}

// GetJudgesByTags 根據標籤獲取題目
func (s *JudgeService) GetJudgesByTags(tags []string) ([]models.Judge, error) {
	if len(tags) == 0 {
		return nil, errors.New("請至少提供一個標籤")
	}

	return s.repo.GetJudgesByTags(tags)
}

// GetJudgesByDifficulty 根據難度獲取題目
func (s *JudgeService) GetJudgesByDifficulty(difficulty int) ([]models.Judge, error) {
	if difficulty < 0 {
		return nil, errors.New("難度值不能為負數")
	}

	return s.repo.GetJudgesByDifficulty(difficulty)
}
