package services

import (
	"fmt"
	"log/slog"

	"github.com/dermot10/code-reviewer/backend_go/models"
	"gorm.io/gorm"
)

type ChatService struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewChatService(db *gorm.DB, logger *slog.Logger) *ChatService {
	return &ChatService{db: db, logger: logger}
}

func (s *ChatService) CreateConversation(userID uint, title string) (*models.Conversation, error) {
	conv := &models.Conversation{
		UserID: userID,
		Title:  title,
	}

	if err := s.db.Create(conv).Error; err != nil {
		return nil, err
	}

	return conv, nil
}

func (s *ChatService) ListConversations(userID uint) ([]models.Conversation, error) {
	var conversations []models.Conversation

	if err := s.db.
		Where("user_id = ? ", userID).
		Where("archived = ?", false).
		Order("updated_at DESC").
		Find(&conversations).Error; err != nil {
		return nil, err
	}
	return conversations, nil
}

func (s *ChatService) CreateMessage(userID, conversationID uint, role, content string) (*models.ChatMessage, error) {

	var conv models.Conversation
	if err := s.db.Where("id = ? AND user_id = ?", conversationID, userID).
		First(&conv).Error; err != nil {
		return nil, fmt.Errorf("conversation not found")
	}

	msg := &models.ChatMessage{
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
	}

	if err := s.db.Create(msg).Error; err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}
	return msg, nil
}

func (s *ChatService) ListMessages(userID, conversationID uint, limit, offset int) ([]models.ChatMessage, error) {

	var conv models.Conversation
	if err := s.db.Where("id = ? AND user_id = ?", conversationID, userID).First(&conv).Error; err != nil {
		return nil, fmt.Errorf("conversation not found")
	}

	var messages []models.ChatMessage

	if err := s.db.Where(
		"conversation_id = ? ", conversationID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *ChatService) ArchiveConversation(userID, conversationID uint) error {
	return s.db.
		Model(&models.Conversation{}).
		Where("id = ? AND user_id = ?", conversationID, userID).
		Update("archived", true).Error
}

func (s *ChatService) RenameConversation(userID, conversationID uint, title string) error {
	return s.db.
		Model(&models.Conversation{}).
		Where("id = ? AND user_id = ?", conversationID, userID).
		Update("title", title).
		Error
}

func (s *ChatService) DeleteConversation(userID, conversationID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", conversationID, userID).
		Delete(&models.Conversation{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("conversation not found or unauthorized")
	}
	return nil
}
