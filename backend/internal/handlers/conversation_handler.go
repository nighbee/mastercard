package handlers

import (
	"strconv"

	"mastercard-backend/internal/services"

	"github.com/gofiber/fiber/v2"
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
func (h *ConversationHandler) GetConversations(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	conversations, total, err := h.conversationService.GetConversations(userID, limit, offset)
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
func (h *ConversationHandler) GetConversation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	conversationID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid conversation ID",
		})
	}

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

