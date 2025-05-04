package repositories

import (
	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type UserProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) *UserProfileRepository {
	return &UserProfileRepository{db: db}
}

func (r *UserProfileRepository) GetAll() ([]models.UserProfile, error) {
	var userProfiles []models.UserProfile
	result := r.db.Preload("User").Preload("School").Preload("Social").Find(&userProfiles)
	return userProfiles, result.Error
}

func (r *UserProfileRepository) GetByID(id uint) (models.UserProfile, error) {
	var userProfile models.UserProfile
	result := r.db.Preload("User").Preload("School").Preload("Social").First(&userProfile, id)
	return userProfile, result.Error
}

func (r *UserProfileRepository) GetByUserID(userID uint) (models.UserProfile, error) {
	var userProfile models.UserProfile
	result := r.db.Where("user_id = ?", userID).
		Preload("User").
		Preload("School").
		Preload("Social").
		First(&userProfile)
	return userProfile, result.Error
}

func (r *UserProfileRepository) Create(userProfile models.UserProfile) (models.UserProfile, error) {
	result := r.db.Create(&userProfile)
	return userProfile, result.Error
}

func (r *UserProfileRepository) Update(userProfile models.UserProfile) error {
	return r.db.Save(&userProfile).Error
}

func (r *UserProfileRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserProfile{}, id).Error
}
