package services

import (
	"errors"
	"time"

	"mastercard-backend/internal/database"
	"mastercard-backend/internal/models"

	"gorm.io/gorm"
)

type ConversationService struct{}

func NewConversationService() *ConversationService {
	return &ConversationService{}
}

// CreateConversation creates a new conversation
func (s *ConversationService) CreateConversation(userID uint, title string) (*models.Conversation, error) {
	conversation := models.Conversation{
		UserID: userID,
		Title:  &title,
	}

	if err := database.DB.Create(&conversation).Error; err != nil {
		return nil, errors.New("failed to create conversation")
	}

	return &conversation, nil
}

// GetConversations retrieves all conversations for a user
func (s *ConversationService) GetConversations(userID uint, limit, offset int) ([]models.Conversation, int64, error) {
	var conversations []models.Conversation
	var total int64

	query := database.DB.Where("user_id = ?", userID)

	// Count total
	if err := query.Model(&models.Conversation{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get conversations
	if err := query.Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&conversations).Error; err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// GetConversation retrieves a single conversation with messages
func (s *ConversationService) GetConversation(conversationID, userID uint) (*models.Conversation, error) {
	var conversation models.Conversation
	if err := database.DB.Where("id = ? AND user_id = ?", conversationID, userID).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		First(&conversation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("conversation not found")
		}
		return nil, err
	}

	return &conversation, nil
}

// UpdateConversation updates a conversation (e.g., rename)
func (s *ConversationService) UpdateConversation(conversationID, userID uint, title string) (*models.Conversation, error) {
	var conversation models.Conversation
	if err := database.DB.Where("id = ? AND user_id = ?", conversationID, userID).
		First(&conversation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("conversation not found")
		}
		return nil, err
	}

	conversation.Title = &title
	conversation.UpdatedAt = time.Now()

	if err := database.DB.Save(&conversation).Error; err != nil {
		return nil, errors.New("failed to update conversation")
	}

	return &conversation, nil
}

// DeleteConversation deletes a conversation
func (s *ConversationService) DeleteConversation(conversationID, userID uint) error {
	result := database.DB.Where("id = ? AND user_id = ?", conversationID, userID).
		Delete(&models.Conversation{})

	if result.Error != nil {
		return errors.New("failed to delete conversation")
	}

	if result.RowsAffected == 0 {
		return errors.New("conversation not found")
	}

	return nil
}

// CreateBranch creates a new conversation branch from a message
func (s *ConversationService) CreateBranch(parentConversationID, branchPointMessageID, userID uint, title string) (*models.Conversation, error) {
	// Verify parent conversation belongs to user
	var parent models.Conversation
	if err := database.DB.Where("id = ? AND user_id = ?", parentConversationID, userID).
		First(&parent).Error; err != nil {
		return nil, errors.New("parent conversation not found")
	}

	// Create new branch
	branch := models.Conversation{
		UserID:             userID,
		Title:              &title,
		ParentBranchID:     &parentConversationID,
		BranchPointMessageID: &branchPointMessageID,
	}

	if err := database.DB.Create(&branch).Error; err != nil {
		return nil, errors.New("failed to create branch")
	}

	return &branch, nil
}

// SearchConversations searches conversations by keyword
func (s *ConversationService) SearchConversations(userID uint, keyword string, limit, offset int) ([]models.Conversation, int64, error) {
	var conversations []models.Conversation
	var total int64

	query := database.DB.Where("user_id = ?", userID).
		Where("to_tsvector('english', COALESCE(title, '')) @@ plainto_tsquery('english', ?)", keyword)

	// Count total
	if err := query.Model(&models.Conversation{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get conversations
	if err := query.Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&conversations).Error; err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

