package repositories

import (
	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type SupportLangRepository struct {
	db *gorm.DB
}

func NewSupportLangRepository(db *gorm.DB) *SupportLangRepository {
	return &SupportLangRepository{db: db}
}

// GetAll 獲取所有支援語言
func (r *SupportLangRepository) GetAll() ([]models.SupportLang, error) {
	var langs []models.SupportLang
	result := r.db.Find(&langs)
	return langs, result.Error
}

// GetByID 根據ID獲取支援語言
func (r *SupportLangRepository) GetByID(id uint) (models.SupportLang, error) {
	var lang models.SupportLang
	result := r.db.First(&lang, id)
	return lang, result.Error
}

// GetByName 根據名稱獲取支援語言
func (r *SupportLangRepository) GetByName(name string) (models.SupportLang, error) {
	var lang models.SupportLang
	result := r.db.Where("name = ?", name).First(&lang)
	return lang, result.Error
}

// GetActiveLanguages 獲取所有啟用的支援語言
func (r *SupportLangRepository) GetActiveLanguages() ([]models.SupportLang, error) {
	var langs []models.SupportLang
	result := r.db.Where("is_active = ?", true).Find(&langs)
	return langs, result.Error
}

// Create 創建新支援語言
func (r *SupportLangRepository) Create(lang models.SupportLang) (models.SupportLang, error) {
	result := r.db.Create(&lang)
	return lang, result.Error
}

// Update 更新支援語言
func (r *SupportLangRepository) Update(lang models.SupportLang) error {
	return r.db.Save(&lang).Error
}

// Delete 刪除支援語言
func (r *SupportLangRepository) Delete(id uint) error {
	return r.db.Delete(&models.SupportLang{}, id).Error
}

// ToggleActive 切換支援語言的啟用狀態
func (r *SupportLangRepository) ToggleActive(id uint, isActive bool) error {
	return r.db.Model(&models.SupportLang{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}
