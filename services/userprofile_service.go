package services

import (
	"errors"
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/repositories"
)

type UserProfileService struct {
	repo     *repositories.UserProfileRepository
	userRepo *repositories.UserRepository
}

func NewUserProfileService(repo *repositories.UserProfileRepository, userRepo *repositories.UserRepository) *UserProfileService {
	return &UserProfileService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *UserProfileService) GetAllUserProfiles() ([]models.UserProfile, error) {
	return s.repo.GetAll()
}

func (s *UserProfileService) GetUserProfileByID(id uint) (models.UserProfile, error) {
	if id == 0 {
		return models.UserProfile{}, errors.New("invalid ID")
	}
	return s.repo.GetByID(id)
}

func (s *UserProfileService) GetUserProfileByUserID(userID uint) (models.UserProfile, error) {
	if userID == 0 {
		return models.UserProfile{}, errors.New("invalid user ID")
	}
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return models.UserProfile{}, errors.New("user not found")
	}

	return s.repo.GetByUserID(userID)
}

func (s *UserProfileService) CreateUserProfile(userProfile models.UserProfile) (models.UserProfile, error) {
	if userProfile.UserID == 0 {
		return models.UserProfile{}, errors.New("user ID is required")
	}

	existingProfile, err := s.repo.GetByUserID(uint(userProfile.UserID))
	if err == nil && existingProfile.ID != 0 {
		return models.UserProfile{}, errors.New("user profile already exists")
	}

	if userProfile.SchoolID == nil || userProfile.SocialID == 0 {
		return models.UserProfile{}, errors.New("school ID and social ID are required")
	}

	return s.repo.Create(userProfile)
}

func (s *UserProfileService) UpdateUserProfile(userProfile models.UserProfile) error {
	return s.repo.Update(userProfile)
}

func (s *UserProfileService) DeleteUserProfile(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("user profile not found")
	}

	return s.repo.Delete(id)
}

func (s *UserProfileService) UpdateAvatarURL(userID uint, avatarURL string) error {
	if avatarURL == "" {
		return errors.New("avatar URL cannot be empty")
	}

	profile, err := s.repo.GetByUserID(userID)
	if err != nil {
		return errors.New("user profile not found")
	}

	profile.AvatarURL = avatarURL
	if err := s.repo.Update(profile); err != nil {
		return errors.New("failed to update avatar URL")
	}
	return nil
}
