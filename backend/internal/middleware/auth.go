package middleware

import (
	"strings"

	"mastercard-backend/internal/database"
	"mastercard-backend/internal/models"
	"mastercard-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token := parts[1]
		claims, err := utils.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Get user from database
		var user models.User
		if err := database.DB.Preload("Role").First(&user, claims.UserID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		if !user.IsActive {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "User account is inactive",
			})
		}

		// Set user context
		c.Locals("user", &user)
		c.Locals("userID", user.ID)
		c.Locals("roleID", user.RoleID)

		// Set RLS context for database queries
		if err := database.SetCurrentUserID(user.ID); err != nil {
			// Log error but don't fail the request
			_ = err
		}

		return c.Next()
	}
}

// OptionalAuthMiddleware allows requests with or without authentication
func OptionalAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				claims, err := utils.ValidateToken(parts[1])
				if err == nil {
					var user models.User
					if err := database.DB.Preload("Role").First(&user, claims.UserID).Error; err == nil && user.IsActive {
						c.Locals("user", &user)
						c.Locals("userID", user.ID)
						c.Locals("roleID", user.RoleID)
						_ = database.SetCurrentUserID(user.ID)
					}
				}
			}
		}
		return c.Next()
	}
}

