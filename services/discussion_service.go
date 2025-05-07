package services

import (
	"errors"
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/repositories"
)

type DiscussionService struct {
	repo     *repositories.DiscussionRepository
	userRepo *repositories.UserRepository
}

func NewDiscussionService(repo *repositories.DiscussionRepository, userRepo *repositories.UserRepository) *DiscussionService {
	return &DiscussionService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *DiscussionService) GetAllDiscussions() ([]models.Discussion, error) {
	return s.repo.GetAll()
}

func (s *DiscussionService) GetDiscussionByID(id uint) (models.Discussion, error) {
	if id == 0 {
		return models.Discussion{}, errors.New("invalid ID")
	}

	// 增加瀏覽量
	s.repo.IncrementViews(id)

	return s.repo.GetByID(id)
}

func (s *DiscussionService) GetUserDiscussions(userID uint) ([]models.Discussion, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.repo.GetByUserID(userID)
}

func (s *DiscussionService) GetChildDiscussions(parentID int, parentType string) ([]models.Discussion, error) {
	if parentID == 0 || parentType == "" {
		return nil, errors.New("parent ID and type are required")
	}

	return s.repo.GetByParent(parentID, parentType)
}

func (s *DiscussionService) CreateDiscussion(discussion models.Discussion) (models.Discussion, error) {

	if discussion.Title == "" {
		return models.Discussion{}, errors.New("title is required")
	}

	if discussion.Content == "" {
		return models.Discussion{}, errors.New("content is required")
	}

	if discussion.Type == "" {
		return models.Discussion{}, errors.New("type is required")
	}

	if discussion.AuthorID == 0 {
		return models.Discussion{}, errors.New("author ID is required")
	}
	_, err := s.userRepo.GetByID(uint(discussion.AuthorID))
	if err != nil {
		return models.Discussion{}, errors.New("author not found")
	}
	return s.repo.Create(discussion)
}

func (s *DiscussionService) UpdateDiscussion(id uint, discussionUpdate models.Discussion) error {

	existingDiscussion, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("discussion not found")
	}

	if discussionUpdate.Title != "" {
		existingDiscussion.Title = discussionUpdate.Title
	}

	if discussionUpdate.Content != "" {
		existingDiscussion.Content = discussionUpdate.Content
	}

	if len(discussionUpdate.Tags) > 0 {
		existingDiscussion.Tags = discussionUpdate.Tags
	}

	return s.repo.Update(existingDiscussion)
}

func (s *DiscussionService) DeleteDiscussion(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("discussion not found")
	}

	return s.repo.Delete(id)
}

func (s *DiscussionService) SearchDiscussions(keyword string) ([]models.Discussion, error) {
	if keyword == "" {
		return nil, errors.New("search keyword is required")
	}

	return s.repo.Search(keyword)
}

func (s *DiscussionService) GetDiscussionsByTags(tags []string) ([]models.Discussion, error) {
	if len(tags) == 0 {
		return nil, errors.New("at least one tag is required")
	}

	return s.repo.GetByTags(tags)
}

func (s *DiscussionService) GetMostViewedDiscussions(limit int) ([]models.Discussion, error) {
	if limit <= 0 {
		limit = 10
	}

	return s.repo.GetMostViewed(limit)
}

func (s *DiscussionService) GetDiscussionsWithFilters(discussionType string, tags []string, keyword string, limit int) ([]models.Discussion, error) {
	return s.repo.GetWithFilters(discussionType, tags, keyword, limit)
}
