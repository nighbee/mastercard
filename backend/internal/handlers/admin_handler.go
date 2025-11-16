package handlers

import (
	"strconv"

	"mastercard-backend/internal/database"
	"mastercard-backend/internal/models"
	"mastercard-backend/internal/middleware"
	"mastercard-backend/internal/services"
	"mastercard-backend/internal/utils"

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

// GetUsers retrieves all users (manager and admin only)
func (h *AdminHandler) GetUsers(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Check if user can manage users (manager or admin)
	if !middleware.CanManageUsers(user) && !middleware.IsAdmin(user) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions to view users",
		})
	}

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

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
	RoleID   *uint  `json:"role_id"`
}

// CreateUser creates a new user (manager and admin only)
func (h *AdminHandler) CreateUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Check if user can manage users (manager or admin)
	if !middleware.CanManageUsers(user) && !middleware.IsAdmin(user) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions to create users",
		})
	}

	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email, password, and full_name are required",
		})
	}

	// Check if email already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User with this email already exists",
		})
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create user
	newUser := models.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	database.DB.Preload("Role").First(&newUser, newUser.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": newUser,
	})
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty"`
	FullName *string `json:"full_name,omitempty"`
	RoleID   *uint   `json:"role_id,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// UpdateUser updates a user (manager and admin only)
func (h *AdminHandler) UpdateUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Check if user can manage users (manager or admin)
	if !middleware.CanManageUsers(user) && !middleware.IsAdmin(user) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions to update users",
		})
	}

	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var targetUser models.User
	if err := database.DB.First(&targetUser, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields if provided
	if req.Email != nil {
		// Check if email already exists for another user
		var existingUser models.User
		if err := database.DB.Where("email = ? AND id != ?", *req.Email, userID).First(&existingUser).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already in use",
			})
		}
		targetUser.Email = *req.Email
	}
	if req.FullName != nil {
		targetUser.FullName = *req.FullName
	}
	if req.RoleID != nil {
		// Verify role exists
		var role models.Role
		if err := database.DB.First(&role, *req.RoleID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid role ID",
			})
		}
		targetUser.RoleID = req.RoleID
	}
	if req.IsActive != nil {
		targetUser.IsActive = *req.IsActive
	}

	if err := database.DB.Save(&targetUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	database.DB.Preload("Role").First(&targetUser, targetUser.ID)

	return c.JSON(fiber.Map{
		"user": targetUser,
	})
}

// DeleteUser deletes a user (admin only)
func (h *AdminHandler) DeleteUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Only admin can delete users
	if !middleware.CanDeleteUsers(user) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions to delete users",
		})
	}

	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Prevent deleting yourself
	if user.ID == uint(userID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete your own account",
		})
	}

	var targetUser models.User
	if err := database.DB.First(&targetUser, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if err := database.DB.Delete(&targetUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// GetAuditLogs retrieves audit logs (manager and admin only)
func (h *AdminHandler) GetAuditLogs(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Check if user can view audit logs (manager or admin)
	if !middleware.HasPermission(user, "audit_logs", "read") && !middleware.IsAdmin(user) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions to view audit logs",
		})
	}
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

// GetMetrics returns system metrics (manager and admin only)
func (h *AdminHandler) GetMetrics(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Check if user can view metrics (manager or admin)
	if !middleware.IsManagerOrHigher(user) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions to view metrics",
		})
	}
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

