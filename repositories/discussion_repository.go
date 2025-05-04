package repositories

import (
	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type DiscussionRepository struct {
	db *gorm.DB
}

func NewDiscussionRepository(db *gorm.DB) *DiscussionRepository {
	return &DiscussionRepository{db: db}
}

func (r *DiscussionRepository) GetAll() ([]models.Discussion, error) {
	var discussions []models.Discussion
	result := r.db.Find(&discussions)
	return discussions, result.Error
}

func (r *DiscussionRepository) GetByID(id uint) (models.Discussion, error) {
	var discussion models.Discussion
	result := r.db.Where("id = ?", id).
		First(&discussion)
	return discussion, result.Error
}

func (r *DiscussionRepository) GetByUserID(userID uint) ([]models.Discussion, error) {
	var discussions []models.Discussion
	result := r.db.Where("author_id = ?", userID).
		Find(&discussions)
	return discussions, result.Error
}

func (r *DiscussionRepository) Create(discussion models.Discussion) (models.Discussion, error) {
	result := r.db.Create(&discussion)
	return discussion, result.Error
}

func (r *DiscussionRepository) Update(discussion models.Discussion) error {
	return r.db.Save(&discussion).Error
}

func (r *DiscussionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Discussion{}, id).Error
}
