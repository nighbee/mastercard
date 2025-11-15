package handlers

import (
	"mastercard-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type QueryHandler struct {
	queryService *services.QueryService
}

func NewQueryHandler(queryService *services.QueryService) *QueryHandler {
	return &QueryHandler{
		queryService: queryService,
	}
}

type QueryRequest struct {
	Query          string `json:"query" validate:"required"`
	ConversationID *uint  `json:"conversation_id,omitempty"`
}

// ExecuteQuery handles natural language query execution
func (h *QueryHandler) ExecuteQuery(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req QueryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query is required",
		})
	}

	message, err := h.queryService.ExecuteQuery(userID, req.Query, req.ConversationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": message,
	})
}

