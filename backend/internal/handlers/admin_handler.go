package handlers

import (
	"strconv"

	"mastercard-backend/internal/database"
	"mastercard-backend/internal/models"
	"mastercard-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	auditService *services.AuditService
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		auditService: services.NewAuditService(),
	}
}

// GetUsers retrieves all users (admin only)
func (h *AdminHandler) GetUsers(c *fiber.Ctx) error {
	var users []models.User
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	if err := database.DB.Preload("Role").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	var total int64
	database.DB.Model(&models.User{}).Count(&total)

	return c.JSON(fiber.Map{
		"users": users,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

// GetAuditLogs retrieves audit logs (admin only)
func (h *AdminHandler) GetAuditLogs(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	action := c.Query("action")
	resource := c.Query("resource")
	status := c.Query("status")

	if limit > 100 {
		limit = 100
	}

	var actionPtr, resourcePtr, statusPtr *string
	if action != "" {
		actionPtr = &action
	}
	if resource != "" {
		resourcePtr = &resource
	}
	if status != "" {
		statusPtr = &status
	}

	logs, total, err := h.auditService.GetAuditLogs(nil, limit, offset, actionPtr, resourcePtr, statusPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"logs":  logs,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

// GetMetrics returns system metrics (admin only)
func (h *AdminHandler) GetMetrics(c *fiber.Ctx) error {
	var totalUsers, activeUsers, totalConversations, totalMessages, totalQueries int64

	database.DB.Model(&models.User{}).Count(&totalUsers)
	database.DB.Model(&models.User{}).Where("is_active = ?", true).Count(&activeUsers)
	database.DB.Model(&models.Conversation{}).Count(&totalConversations)
	database.DB.Model(&models.Message{}).Count(&totalMessages)
	database.DB.Model(&models.AuditLog{}).Where("action = ?", "query").Count(&totalQueries)

	return c.JSON(fiber.Map{
		"metrics": fiber.Map{
			"total_users":        totalUsers,
			"active_users":       activeUsers,
			"total_conversations": totalConversations,
			"total_messages":     totalMessages,
			"total_queries":      totalQueries,
		},
	})
}

