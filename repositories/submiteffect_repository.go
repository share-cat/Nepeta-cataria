package repositories

import (
	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type SubmitEffectRepository struct {
	db *gorm.DB
}

func NewSubmitEffectRepository(db *gorm.DB) *SubmitEffectRepository {
	return &SubmitEffectRepository{db: db}
}

func (r *SubmitEffectRepository) GetAll() ([]models.SubmitEffect, error) {
	var submitEffects []models.SubmitEffect
	result := r.db.Find(&submitEffects)
	return submitEffects, result.Error
}

func (r *SubmitEffectRepository) GetByID(id uint) (models.SubmitEffect, error) {
	var submitEffect models.SubmitEffect
	result := r.db.Where("id = ?", id).First(&submitEffect)
	return submitEffect, result.Error
}

func (r *SubmitEffectRepository) SearchByName(name string) (models.SubmitEffect, error) {
	var submitEffect models.SubmitEffect
	result := r.db.Where("name like ?", name+"%").First(&submitEffect)
	if result.Error != nil {
		return submitEffect, result.Error
	}
	if result.RowsAffected == 0 {
		return submitEffect, gorm.ErrRecordNotFound
	}
	return submitEffect, nil
}

func (r *SubmitEffectRepository) Create(submitEffect models.SubmitEffect) (models.SubmitEffect, error) {
	result := r.db.Create(&submitEffect)
	return submitEffect, result.Error
}

func (r *SubmitEffectRepository) Update(submitEffect models.SubmitEffect) error {
	return r.db.Save(&submitEffect).Error
}

func (r *SubmitEffectRepository) Delete(id uint) error {
	return r.db.Delete(&models.SubmitEffect{}, id).Error
}
