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

func (r *DiscussionRepository) GetByTags(tags []string) ([]models.Discussion, error) {
	var discussions []models.Discussion

	query := r.db
	for _, tag := range tags {
		query = query.Or("? = ANY(tags)", tag)
	}

	result := query.Preload("Author").Find(&discussions)
	return discussions, result.Error
}

func (r *DiscussionRepository) GetWithFilters(discussionType string, tags []string, keyword string, limit int) ([]models.Discussion, error) {
	query := r.db.Model(&models.Discussion{})

	if discussionType != "" {
		query = query.Where("type = ?", discussionType)
	}

	if len(tags) > 0 {
		for _, tag := range tags {
			query = query.Where("? = ANY(tags)", tag)
		}
	}

	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if limit > 0 {
		query = query.Order("views DESC").Limit(limit)
	}

	var discussions []models.Discussion
	result := query.Preload("Author").Find(&discussions)
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
