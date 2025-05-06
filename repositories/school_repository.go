package repositories

import (
	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type SchoolRepository struct {
	db *gorm.DB
}

func NewSchoolRepository(db *gorm.DB) *SchoolRepository {
	return &SchoolRepository{db: db}
}

func (r *SchoolRepository) GetAll() ([]models.School, error) {
	var schools []models.School
	result := r.db.Find(&schools)
	return schools, result.Error
}
func (r *SchoolRepository) GetByID(id uint) (models.School, error) {
	var school models.School
	result := r.db.First(&school, id)
	return school, result.Error
}

func (r *SchoolRepository) SearchByName(name string) (models.School, error) {
	var school models.School
	result := r.db.Where("school_name like ?", name+"%").First(&school)
	return school, result.Error
}

func (r *SchoolRepository) Create(school models.School) (models.School, error) {
	result := r.db.Create(&school)
	return school, result.Error
}

func (r *SchoolRepository) Update(school models.School) error {
	return r.db.Save(&school).Error
}

func (r *SchoolRepository) Delete(id uint) error {
	return r.db.Delete(&models.School{}, id).Error
}
