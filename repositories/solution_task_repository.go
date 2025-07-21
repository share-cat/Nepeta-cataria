package repositories

import (
	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type SolutionTaskRepository struct {
	db *gorm.DB
}

func NewSolutionTaskRepository(db *gorm.DB) *SolutionTaskRepository {
	return &SolutionTaskRepository{db: db}
}

// Create 創建新的解答任務
func (r *SolutionTaskRepository) Create(task models.SolutionTask) (models.SolutionTask, error) {
	err := r.db.Create(&task).Error
	return task, err
}

// GetByID 根據ID獲取解答任務
func (r *SolutionTaskRepository) GetByID(id uint) (models.SolutionTask, error) {
	var task models.SolutionTask
	err := r.db.First(&task, id).Error
	return task, err
}

// GetBySolutionID 根據Solution ID獲取所有任務
func (r *SolutionTaskRepository) GetBySolutionID(solutionID uint) ([]models.SolutionTask, error) {
	var tasks []models.SolutionTask
	err := r.db.Where("solution_id = ?", solutionID).Order("task_index").Find(&tasks).Error
	return tasks, err
}

// Update 更新解答任務
func (r *SolutionTaskRepository) Update(task models.SolutionTask) error {
	return r.db.Save(&task).Error
}

// Delete 刪除解答任務
func (r *SolutionTaskRepository) Delete(id uint) error {
	return r.db.Delete(&models.SolutionTask{}, id).Error
}

// DeleteBySolutionID 根據Solution ID刪除所有相關任務
func (r *SolutionTaskRepository) DeleteBySolutionID(solutionID uint) error {
	return r.db.Where("solution_id = ?", solutionID).Delete(&models.SolutionTask{}).Error
}

// UpdateStatus 更新任務狀態
func (r *SolutionTaskRepository) UpdateStatus(id uint, status string, executionTime, memoryUsage int64, score int) error {
	return r.db.Model(&models.SolutionTask{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         status,
			"execution_time": executionTime,
			"memory_usage":   memoryUsage,
			"score":          score,
		}).Error
}

// CreateBatch 批量創建任務
func (r *SolutionTaskRepository) CreateBatch(tasks []models.SolutionTask) error {
	return r.db.CreateInBatches(tasks, 100).Error
}
