package repositories

import (
	"time"

	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type SolutionRepository struct {
	db *gorm.DB
}

func NewSolutionRepository(db *gorm.DB) *SolutionRepository {
	return &SolutionRepository{db: db}
}

// GetAll 獲取所有答案
func (r *SolutionRepository) GetAll() ([]models.Solution, error) {
	var solutions []models.Solution
	result := r.db.Select("solutions.*").
		Preload("User").
		Preload("Judge").
		Preload("Language").
		Preload("Tasks").
		Find(&solutions)
	return solutions, result.Error
}

// GetByID 根據ID獲取答案
func (r *SolutionRepository) GetByID(id uint) (models.Solution, error) {
	var solution models.Solution
	result := r.db.Select("solutions.*").
		Preload("User").
		Preload("Judge").
		Preload("Language").
		Preload("Tasks").
		First(&solution, id)
	return solution, result.Error
}

// GetByUserID 獲取特定用戶的所有答案
func (r *SolutionRepository) GetByUserID(userID uint) ([]models.Solution, error) {
	var solutions []models.Solution
	result := r.db.Where("user_id = ?", userID).
		Preload("Judge").
		Preload("Language").
		Preload("Tasks").
		Find(&solutions)
	return solutions, result.Error
}

// GetByJudgeID 獲取特定題目的所有答案
func (r *SolutionRepository) GetByJudgeID(judgeID uint) ([]models.Solution, error) {
	var solutions []models.Solution
	result := r.db.Where("judge_id = ?", judgeID).
		Preload("User").
		Preload("Language").
		Preload("Tasks").
		Find(&solutions)
	return solutions, result.Error
}

// Create 創建新答案
func (r *SolutionRepository) Create(solution models.Solution) (models.Solution, error) {
	solution.SubmittedAt = time.Now()
	solution.Status = "pending" // 默認狀態

	result := r.db.Create(&solution)
	return solution, result.Error
}

// Update 更新答案
func (r *SolutionRepository) Update(solution models.Solution) error {
	return r.db.Save(&solution).Error
}

// UpdateStatus 更新答案狀態 (已棄用，改用 UpdateSolutionStatus)
func (r *SolutionRepository) UpdateStatus(id uint, status string, executionTime, memoryUsage, score int) error {
	return r.db.Model(&models.Solution{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": status,
		}).Error
}

// UpdateSolutionStatus 更新解答狀態
func (r *SolutionRepository) UpdateSolutionStatus(id uint, status string) error {
	return r.db.Model(&models.Solution{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 刪除答案
func (r *SolutionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Solution{}, id).Error
}

// GetLatestSolutions 獲取最新提交的答案
func (r *SolutionRepository) GetLatestSolutions(limit int) ([]models.Solution, error) {
	var solutions []models.Solution
	result := r.db.Order("submitted_at desc").
		Limit(limit).
		Preload("User").
		Preload("Judge").
		Preload("Language").
		Preload("Tasks").
		Find(&solutions)
	return solutions, result.Error
}

// GetBestSolutions 獲取評分最高的答案
func (r *SolutionRepository) GetBestSolutions(judgeID uint, limit int) ([]models.Solution, error) {
	var solutions []models.Solution

	subQuery := r.db.Table("solution_tasks").
		Select("solution_id, SUM(execution_time) as total_execution_time, SUM(memory_usage) as total_memory_usage").
		Where("status = ?", "accepted").
		Group("solution_id")

	query := r.db.Table("solutions").
		Select("solutions.*, task_stats.total_execution_time, task_stats.total_memory_usage").
		Joins("JOIN (?) as task_stats ON solutions.id = task_stats.solution_id", subQuery).
		Where("solutions.status = ?", "accepted")

	if judgeID > 0 {
		query = query.Where("solutions.judge_id = ?", judgeID)
	}

	result := query.Order("task_stats.total_memory_usage ASC, task_stats.total_execution_time ASC").
		Limit(limit).
		Preload("User").
		Preload("Judge").
		Preload("Language").
		Preload("Tasks").
		Find(&solutions)
	return solutions, result.Error
}
