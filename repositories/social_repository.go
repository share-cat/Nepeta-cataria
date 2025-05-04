package repositories

import (
	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type SociaRepository struct {
	db *gorm.DB
}

func NewSociaRepository(db *gorm.DB) *SociaRepository {
	return &SociaRepository{db: db}
}

func (r *SociaRepository) GetAll() ([]models.Social, error) {
	var social []models.Social
	result := r.db.Find(&social)
	return social, result.Error
}

func (r *SociaRepository) GetByID(id uint) (models.Social, error) {
	var social models.Social
	result := r.db.First(&social, id)
	return social, result.Error
}
func (r *SociaRepository) Create(social models.Social) (models.Social, error) {
	result := r.db.Create(&social)
	return social, result.Error
}

func (r *SociaRepository) Update(social models.Social) error {
	return r.db.Save(&social).Error
}

func (r *SociaRepository) Delete(id uint) error {
	return r.db.Delete(&models.Social{}, id).Error
}
