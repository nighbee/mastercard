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

// HasRole checks if user has one of the specified roles
func HasRole(user *models.User, roleNames ...string) bool {
	if user == nil || user.RoleID == nil {
		return false
	}

	var role models.Role
	if err := database.DB.First(&role, *user.RoleID).Error; err != nil {
		return false
	}

	for _, roleName := range roleNames {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// HasPermission checks if user has a specific permission
func HasPermission(user *models.User, resource, action string) bool {
	if user == nil || user.RoleID == nil {
		return false
	}

	var permission models.Permission
	if err := database.DB.Where("resource = ? AND action = ?", resource, action).First(&permission).Error; err != nil {
		return false
	}

	var count int64
	if err := database.DB.Table("role_permissions").
		Where("role_id = ? AND permission_id = ?", *user.RoleID, permission.ID).
		Count(&count).Error; err != nil || count == 0 {
		return false
	}

	return true
}

// HasRoleOrHigher checks if user has a role that is equal to or higher than the required role
// Hierarchy: admin > manager > analyzer
func HasRoleOrHigher(user *models.User, requiredRole string) bool {
	if user == nil || user.RoleID == nil {
		return false
	}

	var role models.Role
	if err := database.DB.First(&role, *user.RoleID).Error; err != nil {
		return false
	}

	// Define role hierarchy (higher number = more privileges)
	roleHierarchy := map[string]int{
		"analyzer": 1,
		"manager":  2,
		"admin":    3,
	}

	userLevel, userExists := roleHierarchy[role.Name]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevel
}

// IsAdmin checks if user is admin
func IsAdmin(user *models.User) bool {
	return HasRole(user, "admin")
}

// IsManagerOrHigher checks if user is manager or admin
func IsManagerOrHigher(user *models.User) bool {
	return HasRoleOrHigher(user, "manager")
}

// CanManageUsers checks if user can manage users (manager or admin)
func CanManageUsers(user *models.User) bool {
	return HasPermission(user, "users", "read") && HasPermission(user, "users", "create")
}

// CanDeleteUsers checks if user can delete users (admin only)
func CanDeleteUsers(user *models.User) bool {
	return HasPermission(user, "users", "delete")
}

// CanViewAllConversations checks if user can view all conversations (manager or admin)
func CanViewAllConversations(user *models.User) bool {
	return HasPermission(user, "conversations", "read_all")
}

// CanDeleteAllConversations checks if user can delete any conversation (admin only)
func CanDeleteAllConversations(user *models.User) bool {
	return HasPermission(user, "conversations", "delete_all")
}

// CanConfigureSystem checks if user can configure system (admin only)
func CanConfigureSystem(user *models.User) bool {
	return HasPermission(user, "system", "configure")
}
