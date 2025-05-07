package services

import (
	"errors"
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/repositories"
)

type ChatSessionService struct {
	repo      *repositories.ChatSessionRepository
	userRepo  *repositories.UserRepository
	judgeRepo *repositories.JudgeRepository
}

func NewChatSessionService(repo *repositories.ChatSessionRepository, userRepo *repositories.UserRepository, judgeRepo *repositories.JudgeRepository) *ChatSessionService {
	return &ChatSessionService{
		repo:      repo,
		userRepo:  userRepo,
		judgeRepo: judgeRepo,
	}
}

func (s *ChatSessionService) GetAllChatSessions() ([]models.ChatSession, error) {
	return s.repo.GetAll()
}

func (s *ChatSessionService) GetChatSessionByID(id uint) (models.ChatSession, error) {
	if id == 0 {
		return models.ChatSession{}, errors.New("invalid ID")
	}
	return s.repo.GetByID(id)
}

func (s *ChatSessionService) GetUserChatSessions(userID uint) ([]models.ChatSession, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.repo.GetByUserID(userID)
}

func (s *ChatSessionService) GetJudgeChatSessions(judgeID uint) ([]models.ChatSession, error) {
	if judgeID == 0 {
		return nil, errors.New("invalid judge ID")
	}

	_, err := s.judgeRepo.GetByID(judgeID)
	if err != nil {
		return nil, errors.New("judge not found")
	}

	return s.repo.GetByJudgeID(judgeID)
}

func (s *ChatSessionService) CreateChatSession(chatSession models.ChatSession) (models.ChatSession, error) {
	if chatSession.UserID == 0 || chatSession.JudgeID == 0 {
		return models.ChatSession{}, errors.New("user ID and judge ID are required")
	}

	if chatSession.Title == "" {
		return models.ChatSession{}, errors.New("title is required")
	}

	_, err := s.userRepo.GetByID(uint(chatSession.UserID))
	if err != nil {
		return models.ChatSession{}, errors.New("user not found")
	}

	_, err = s.judgeRepo.GetByID(uint(chatSession.JudgeID))
	if err != nil {
		return models.ChatSession{}, errors.New("judge not found")
	}

	return s.repo.Create(chatSession)
}

func (s *ChatSessionService) UpdateChatSession(id uint, chatSession models.ChatSession) error {

	existingSession, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("chat session not found")
	}

	if chatSession.Title != "" {
		existingSession.Title = chatSession.Title
	}

	return s.repo.Update(existingSession)
}

func (s *ChatSessionService) DeleteChatSession(id uint) error {

	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("chat session not found")
	}

	return s.repo.Delete(id)
}

func (s *ChatSessionService) UpdateChatHistory(id uint, history interface{}) error {

	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("chat session not found")
	}

	return s.repo.UpdateChatHistory(id, history)
}

func (s *ChatSessionService) GetUserJudgeChatSessions(userID, judgeID uint) ([]models.ChatSession, error) {
	if userID == 0 || judgeID == 0 {
		return nil, errors.New("invalid user ID or judge ID")
	}

	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	_, err = s.judgeRepo.GetByID(judgeID)
	if err != nil {
		return nil, errors.New("judge not found")
	}

	return s.repo.GetByUserAndJudge(userID, judgeID)
}

func (s *ChatSessionService) UpdateTitle(id uint, title string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("chat session not found")
	}

	return s.repo.UpdateTitle(id, title)
}
