package repositories

import (
	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type JudgeRepository struct {
	db *gorm.DB
}

func NewJudgeRepository(db *gorm.DB) *JudgeRepository {
	return &JudgeRepository{db: db}
}

// GetAll 獲取所有題目
func (r *JudgeRepository) GetAll() ([]models.Judge, error) {
	var judges []models.Judge
	result := r.db.Preload("Creator").Find(&judges)
	return judges, result.Error
}

// GetJudgesWithFilters 使用篩選條件獲取題目
func (r *JudgeRepository) GetJudgesWithFilters(search, judgeType string, page int, activeOnly bool) ([]models.Judge, int, error) {
	const pageSize = 20
	offset := (page - 1) * pageSize

	query := r.db.Model(&models.Judge{})

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if judgeType != "" {
		query = query.Where("type = ?", judgeType)
	}

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var judges []models.Judge
	if err := query.Offset(offset).Limit(pageSize).Find(&judges).Error; err != nil {
		return nil, 0, err
	}

	return judges, int(total), nil
}

// GetByID 根據ID獲取題目
func (r *JudgeRepository) GetByID(id uint) (models.Judge, error) {
	var judge models.Judge
	result := r.db.Preload("Creator").First(&judge, id)
	return judge, result.Error
}

// GetActiveJudges 獲取所有活躍的題目
func (r *JudgeRepository) GetActiveJudges() ([]models.Judge, error) {
	var judges []models.Judge
	result := r.db.Where("is_active = ?", true).Preload("Creator").Find(&judges)
	return judges, result.Error
}

// GetJudgesByCreator 獲取特定用戶創建的題目
func (r *JudgeRepository) GetJudgesByCreator(creatorID uint) ([]models.Judge, error) {
	var judges []models.Judge
	result := r.db.Where("creator_id = ?", creatorID).Preload("Creator").Find(&judges)
	return judges, result.Error
}

// Create 創建新題目
func (r *JudgeRepository) Create(judge models.Judge) (models.Judge, error) {
	result := r.db.Create(&judge)
	return judge, result.Error
}

// Update 更新題目
func (r *JudgeRepository) Update(judge models.Judge) error {
	return r.db.Save(&judge).Error
}

// Delete 刪除題目
func (r *JudgeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Judge{}, id).Error
}

// GetJudgesByTags 根據標籤獲取題目
func (r *JudgeRepository) GetJudgesByTags(tags []string) ([]models.Judge, error) {
	var judges []models.Judge
	query := r.db.Where("? = ANY(tag)", tags[0])

	for i := 1; i < len(tags); i++ {
		query = query.Or("? = ANY(tag)", tags[i])
	}

	result := query.Preload("Creator").Find(&judges)
	return judges, result.Error
}

// GetJudgesByDifficulty 根據難度獲取題目
func (r *JudgeRepository) GetJudgesByDifficulty(difficulty int) ([]models.Judge, error) {
	var judges []models.Judge
	result := r.db.Where("difficulty = ?", difficulty).Preload("Creator").Find(&judges)
	return judges, result.Error
}
