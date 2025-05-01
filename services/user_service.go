// services/user_service.go
package services

import (
	"errors"
	"time"

	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetAllUsers 獲取所有用戶
func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

// GetUserByID 根據ID獲取用戶
func (s *UserService) GetUserByID(id uint) (models.User, error) {
	return s.repo.GetByID(id)
}

// CreateUser 創建新用戶
// TODO: 權限更新
func (s *UserService) CreateUser(user models.User) (models.User, error) {

	existingUser, err := s.repo.GetByUsername(user.Username)
	if err == nil && existingUser.ID != 0 {
		return models.User{}, errors.New("用戶名已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user.Password = string(hashedPassword)
	user.RegisterTime = time.Now()

	return s.repo.Create(user)
}

// UpdateUser 更新用戶信息
func (s *UserService) UpdateUser(id uint, userUpdate models.User) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if userUpdate.Username != "" {
		if user.Username != userUpdate.Username {
			existingUser, err := s.repo.GetByUsername(userUpdate.Username)
			if err == nil && existingUser.ID != 0 && existingUser.ID != id {
				return errors.New("用戶名已存在")
			}
			user.Username = userUpdate.Username
		}
	}

	if userUpdate.Email != "" {
		user.Email = userUpdate.Email
	}

	if userUpdate.AvatarUrl != "" {
		user.AvatarUrl = userUpdate.AvatarUrl
	}

	if userUpdate.PermissionID != nil {
		user.PermissionID = userUpdate.PermissionID
	}

	return s.repo.Update(user)
}

// UpdateUserPassword 更新用戶密碼
func (s *UserService) UpdateUserPassword(id uint, oldPassword, newPassword string) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("密碼不正確")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.repo.Update(user)
}

// DeleteUser 刪除用戶
func (s *UserService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}

// GetUserProfile 獲取用戶資料
func (s *UserService) GetUserProfile(userID uint) (models.UserProfile, error) {
	return s.repo.GetUserProfile(userID)
}
