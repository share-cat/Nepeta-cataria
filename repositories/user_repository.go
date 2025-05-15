package repositories

import (
	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAll 獲取所有用戶
func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Preload("Permission").Find(&users)
	return users, result.Error
}

// GetByID 根據ID獲取用戶
func (r *UserRepository) GetByID(id uint) (models.User, error) {
	var user models.User
	result := r.db.Preload("Permission").First(&user, id)
	return user, result.Error
}

// GetByUsername 根據用戶名獲取用戶
func (r *UserRepository) GetByUsername(username string) (models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).Preload("Permission").First(&user)
	return user, result.Error
}

// RegisterUser 註冊用戶及相關資料
func (r *UserRepository) RegisterUser(newUser models.User) (models.User, error) {

	tx := r.db.Begin()

	permission := models.Permission{
		IsAdmin:   false,
		IsTeacher: false,
		IsStudent: true,
	}

	if err := tx.Create(&permission).Error; err != nil {
		tx.Rollback()
		return models.User{}, err
	}

	newUser.PermissionID = int(permission.ID)

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return models.User{}, err
	}

	social := models.Social{}
	if err := tx.Create(&social).Error; err != nil {
		tx.Rollback()
		return models.User{}, err
	}

	userProfile := models.UserProfile{
		UserID:    int(newUser.ID),
		SocialID:  int(social.ID),
		AvatarURL: newUser.AvatarUrl,
	}

	if err := tx.Create(&userProfile).Error; err != nil {
		tx.Rollback()
		return models.User{}, err
	}

	rpgProfile := models.RpgProfile{
		UserID: int(newUser.ID),
		Level:  1,
	}

	if err := tx.Create(&rpgProfile).Error; err != nil {
		tx.Rollback()
		return models.User{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return models.User{}, err
	}

	var createdUser models.User
	r.db.Preload("Permission").First(&createdUser, newUser.ID)
	return createdUser, nil
}

// Create 創建新用戶
func (r *UserRepository) Create(user models.User) (models.User, error) {
	result := r.db.Create(&user)
	return user, result.Error
}

// Update 更新用戶
func (r *UserRepository) Update(user models.User) error {
	return r.db.Save(&user).Error
}

// Delete 刪除用戶
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// GetUserProfile 獲取用戶資料
func (r *UserRepository) GetUserProfile(userID uint) (models.UserProfile, error) {
	var profile models.UserProfile
	result := r.db.Where("user_id = ?", userID).
		Preload("School").
		Preload("Social").
		First(&profile)
	return profile, result.Error
}
