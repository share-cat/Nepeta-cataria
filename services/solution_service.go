package services

import (
	"errors"
	"time"

	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/pkg/catwheel"
	"share-cat/Nepeta-cataria/repositories"
)

type SolutionService struct {
	repo      *repositories.SolutionRepository
	userRepo  *repositories.UserRepository
	judgeRepo *repositories.JudgeRepository
	langRepo  *repositories.SupportLangRepository
}

func NewSolutionService(
	repo *repositories.SolutionRepository,
	userRepo *repositories.UserRepository,
	judgeRepo *repositories.JudgeRepository,
	langRepo *repositories.SupportLangRepository,
) *SolutionService {
	return &SolutionService{
		repo:      repo,
		userRepo:  userRepo,
		judgeRepo: judgeRepo,
		langRepo:  langRepo,
	}
}

// GetAllSolutions 獲取所有答案
func (s *SolutionService) GetAllSolutions() ([]models.Solution, error) {
	return s.repo.GetAll()
}

// GetSolutionByID 根據ID獲取答案
func (s *SolutionService) GetSolutionByID(id uint) (models.Solution, error) {
	return s.repo.GetByID(id)
}

// GetUserSolutions 獲取特定用戶的所有答案
func (s *SolutionService) GetUserSolutions(userID uint) ([]models.Solution, error) {
	// 檢查用戶是否存在
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用戶不存在")
	}

	return s.repo.GetByUserID(userID)
}

// GetJudgeSolutions 獲取特定題目的所有答案
func (s *SolutionService) GetJudgeSolutions(judgeID uint) ([]models.Solution, error) {

	_, err := s.judgeRepo.GetByID(judgeID)
	if err != nil {
		return nil, errors.New("題目不存在")
	}

	return s.repo.GetByJudgeID(judgeID)
}

// SubmitSolution 提交新答案
func (s *SolutionService) SubmitSolution(solution models.Solution) (models.Solution, error) {
	// 檢查用戶是否存在
	_, err := s.userRepo.GetByID(uint(solution.UserID))
	if err != nil {
		return models.Solution{}, errors.New("用戶不存在")
	}

	judge, err := s.judgeRepo.GetByID(uint(solution.JudgeID))
	if err != nil {
		return models.Solution{}, errors.New("題目不存在")
	}

	if !judge.IsActive {
		return models.Solution{}, errors.New("題目已停用")
	}

	lang, err := s.langRepo.GetByID(uint(solution.LanguageID))
	if err != nil {
		return models.Solution{}, errors.New("不支援該語言")
	}

	if !lang.IsActive {
		return models.Solution{}, errors.New("該語言已停用")
	}

	solution.SubmittedAt = time.Now()

	solution.Status = "pending"

	solution.ExecutionTime = 0
	solution.MemoryUsage = 0
	solution.Score = 0

	return s.repo.Create(solution)
}

func (s *SolutionService) UpdateSolutionFromExecution(result catwheel.ExecuteResponse, solutionID uint) (models.Solution, error) {
	solution, err := s.repo.GetByID(solutionID)
	if err != nil {
		return models.Solution{}, errors.New("答案不存在")
	}

	solution.Status = result.Result.Status
	solution.ExecutionTime = int(result.Result.ExecutionTime)
	solution.MemoryUsage = int(result.Result.MemoryUsed)

	err = s.repo.Update(solution)
	if err != nil {
		return models.Solution{}, err
	}

	return solution, nil
}

// UpdateSolutionResult 更新答案評分結果
func (s *SolutionService) UpdateSolutionResult(id uint, status string, executionTime, memoryUsage, score int) error {

	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("答案不存在")
	}

	validStatus := map[string]bool{
		"pending":               true,
		"running":               true,
		"accepted":              true,
		"wrong_answer":          true,
		"time_limit_exceeded":   true,
		"memory_limit_exceeded": true,
		"compilation_error":     true,
		"runtime_error":         true,
	}

	if !validStatus[status] {
		return errors.New("無效的狀態")
	}

	return s.repo.UpdateStatus(id, status, executionTime, memoryUsage, score)
}

// DeleteSolution 刪除答案
func (s *SolutionService) DeleteSolution(id uint) error {

	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("答案不存在")
	}

	return s.repo.Delete(id)
}

// GetLatestSolutions 獲取最新提交的答案
func (s *SolutionService) GetLatestSolutions(limit int) ([]models.Solution, error) {
	if limit <= 0 {
		limit = 10
	}

	return s.repo.GetLatestSolutions(limit)
}

// GetBestSolutions 獲取評分最高的答案
func (s *SolutionService) GetBestSolutions(judgeID uint, limit int) ([]models.Solution, error) {
	if limit <= 0 {
		limit = 10
	}

	if judgeID > 0 {
		_, err := s.judgeRepo.GetByID(judgeID)
		if err != nil {
			return nil, errors.New("題目不存在")
		}
	}

	return s.repo.GetBestSolutions(judgeID, limit)
}
