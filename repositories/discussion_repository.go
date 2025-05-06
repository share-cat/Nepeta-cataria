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

func (r *DiscussionRepository) GetByParent(parentID int, parentType string) ([]models.Discussion, error) {
	var discussions []models.Discussion
	result := r.db.Where("pid = ? AND ptype = ?", parentID, parentType).
		Preload("Author").
		Find(&discussions)
	return discussions, result.Error
}

func (r *DiscussionRepository) GetMostViewed(limit int) ([]models.Discussion, error) {
	var discussions []models.Discussion
	result := r.db.Order("views DESC").
		Limit(limit).
		Preload("Author").
		Find(&discussions)
	return discussions, result.Error
}

func (r *DiscussionRepository) Search(keyword string) ([]models.Discussion, error) {
	var discussions []models.Discussion
	result := r.db.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Preload("Author").
		Find(&discussions)
	return discussions, result.Error
}

func (r *DiscussionRepository) IncrementViews(id uint) error {
	return r.db.Model(&models.Discussion{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
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
