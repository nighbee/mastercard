package handlers

import (
	"strconv"

	"mastercard-backend/internal/database"
	"mastercard-backend/internal/middleware"
	"mastercard-backend/internal/models"
	"mastercard-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ConversationHandler struct {
	conversationService *services.ConversationService
}

func NewConversationHandler() *ConversationHandler {
	return &ConversationHandler{
		conversationService: services.NewConversationService(),
	}
}

type CreateConversationRequest struct {
	Title string `json:"title"`
}

type UpdateConversationRequest struct {
	Title string `json:"title" validate:"required"`
}

type CreateBranchRequest struct {
	Title              string `json:"title" validate:"required"`
	BranchPointMessageID uint  `json:"branch_point_message_id" validate:"required"`
}

// CreateConversation creates a new conversation
func (h *ConversationHandler) CreateConversation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req CreateConversationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	conversation, err := h.conversationService.CreateConversation(userID, req.Title)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"conversation": conversation,
	})
}

// GetConversations retrieves all conversations for the user
// Analyzers see only their own conversations, Managers and Admins see all conversations
func (h *ConversationHandler) GetConversations(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	user, ok := c.Locals("user").(*models.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	// Check if user can view all conversations (manager or admin)
	var conversations []models.Conversation
	var total int64
	var err error

	if middleware.CanViewAllConversations(user) {
		// Managers and admins can see all conversations
		conversations, total, err = h.conversationService.GetAllConversations(limit, offset)
	} else {
		// Analyzers see only their own conversations
		conversations, total, err = h.conversationService.GetConversations(userID, limit, offset)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"conversations": conversations,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

// GetConversation retrieves a single conversation with messages
// Analyzers can only view their own conversations, Managers and Admins can view all
func (h *ConversationHandler) GetConversation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	user, ok := c.Locals("user").(*models.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	conversationID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid conversation ID",
		})
	}

	// Check if user can view all conversations (manager or admin)
	if middleware.CanViewAllConversations(user) {
		// Managers and admins can view any conversation
		var conversation models.Conversation
		if err := database.DB.Where("id = ?", conversationID).
			Preload("Messages", func(db *gorm.DB) *gorm.DB {
				return db.Order("created_at ASC")
			}).
			Preload("User").
			First(&conversation).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Conversation not found",
			})
		}
		return c.JSON(fiber.Map{
			"conversation": conversation,
		})
	} else {
		// Analyzers can only view their own conversations
		conversation, err := h.conversationService.GetConversation(uint(conversationID), userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"conversation": conversation,
		})
	}
}

// UpdateConversation updates a conversation
func (h *ConversationHandler) UpdateConversation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	conversationID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid conversation ID",
		})
	}

	var req UpdateConversationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	conversation, err := h.conversationService.UpdateConversation(uint(conversationID), userID, req.Title)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"conversation": conversation,
	})
}

// DeleteConversation deletes a conversation
func (h *ConversationHandler) DeleteConversation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	conversationID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid conversation ID",
		})
	}

	if err := h.conversationService.DeleteConversation(uint(conversationID), userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Conversation deleted successfully",
	})
}

// CreateBranch creates a new conversation branch
func (h *ConversationHandler) CreateBranch(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	parentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid conversation ID",
		})
	}

	var req CreateBranchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	branch, err := h.conversationService.CreateBranch(uint(parentID), req.BranchPointMessageID, userID, req.Title)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"conversation": branch,
	})
}

// SearchConversations searches conversations by keyword
func (h *ConversationHandler) SearchConversations(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	keyword := c.Query("q")
	if keyword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search keyword is required",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	conversations, total, err := h.conversationService.SearchConversations(userID, keyword, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"conversations": conversations,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

