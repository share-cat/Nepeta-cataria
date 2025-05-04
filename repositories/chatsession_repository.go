package repositories

import (
	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type ChatSessionRepository struct {
	db *gorm.DB
}

func NewChatSessionRepository(db *gorm.DB) *ChatSessionRepository {
	return &ChatSessionRepository{db: db}
}

func (r *ChatSessionRepository) GetAll() ([]models.ChatSession, error) {
	var chatSessions []models.ChatSession
	result := r.db.Find(&chatSessions)
	return chatSessions, result.Error
}

func (r *ChatSessionRepository) GetByID(id uint) (models.ChatSession, error) {
	var chatSessions models.ChatSession
	result := r.db.Where("id = ?", id).
		First(&chatSessions)
	return chatSessions, result.Error
}

func (r *ChatSessionRepository) GetByUserID(userID uint) (models.ChatSession, error) {
	var chatSessions models.ChatSession
	result := r.db.Where("user_id = ?", userID).
		First(&chatSessions)
	return chatSessions, result.Error
}

func (r *ChatSessionRepository) GetByJudgeID(judgeID uint) ([]models.ChatSession, error) {
	var chatSessions []models.ChatSession
	result := r.db.Where("judge_id = ?", judgeID).
		Find(&chatSessions)
	return chatSessions, result.Error
}

func (r *ChatSessionRepository) Create(chatSession models.ChatSession) (models.ChatSession, error) {
	result := r.db.Create(&chatSession)
	return chatSession, result.Error
}

func (r *ChatSessionRepository) Update(chatSession models.ChatSession) error {
	return r.db.Save(&chatSession).Error
}

func (r *ChatSessionRepository) Delete(id uint) error {
	return r.db.Delete(&models.ChatSession{}, id).Error
}
