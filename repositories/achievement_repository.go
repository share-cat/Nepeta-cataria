package repositories

import (
	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type AchievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

func (r *AchievementRepository) GetAll() ([]models.Achievement, error) {
	var achievements []models.Achievement
	result := r.db.Find(&achievements)
	return achievements, result.Error
}

func (r *AchievementRepository) GetByID(id uint) (models.Achievement, error) {
	var achievement models.Achievement
	result := r.db.Where("id = ?", id).
		First(&achievement)
	return achievement, result.Error
}

func (r *AchievementRepository) Create(achievement models.Achievement) (models.Achievement, error) {
	result := r.db.Create(&achievement)
	return achievement, result.Error
}

func (r *AchievementRepository) Update(achievement models.Achievement) error {
	return r.db.Save(&achievement).Error
}

func (r *AchievementRepository) Delete(id uint) error {
	return r.db.Delete(&models.Achievement{}, id).Error
}
