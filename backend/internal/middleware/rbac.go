package middleware

import (
	"mastercard-backend/internal/database"
	"mastercard-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

// RequireRole checks if user has one of the required roles
func RequireRole(roleNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.User)
		if !ok || user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		if user.RoleID == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "User has no role assigned",
			})
		}

		var role models.Role
		if err := database.DB.First(&role, *user.RoleID).Error; err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Role not found",
			})
		}

		for _, requiredRole := range roleNames {
			if role.Name == requiredRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}

// RequirePermission checks if user has permission for a resource and action
func RequirePermission(resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.User)
		if !ok || user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		if user.RoleID == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "User has no role assigned",
			})
		}

		var permission models.Permission
		if err := database.DB.Where("resource = ? AND action = ?", resource, action).First(&permission).Error; err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Permission not found",
			})
		}

		// Check if user's role has this permission
		var count int64
		if err := database.DB.Table("role_permissions").
			Where("role_id = ? AND permission_id = ?", *user.RoleID, permission.ID).
			Count(&count).Error; err != nil || count == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

