package repositories

import (
	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type SocialRepository struct {
	db *gorm.DB
}

func NewSociaRepository(db *gorm.DB) *SocialRepository {
	return &SocialRepository{db: db}
}

func (r *SocialRepository) GetAll() ([]models.Social, error) {
	var social []models.Social
	result := r.db.Find(&social)
	return social, result.Error
}

func (r *SocialRepository) GetByID(id uint) (models.Social, error) {
	var social models.Social
	result := r.db.First(&social, id)
	return social, result.Error
}
func (r *SocialRepository) Create(social models.Social) (models.Social, error) {
	result := r.db.Create(&social)
	return social, result.Error
}

func (r *SocialRepository) Update(social models.Social) error {
	return r.db.Save(&social).Error
}

func (r *SocialRepository) Delete(id uint) error {
	return r.db.Delete(&models.Social{}, id).Error
}
