package middleware

import (
	jwtPkg "go-trial/pkg/jwt"
	"go-trial/pkg/response"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RequirePermission is a middleware to check if the authenticated user has a specific permission.
func RequirePermission(db *gorm.DB, requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*jwtPkg.Claims)
		if !ok || claims == nil {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		// Bypass check untuk role programmer dan administrator (akses ke semua route)
		if claims.Role == "programmer" || claims.Role == "administrator" {
			return c.Next()
		}

		var count int64
		err := db.Table("users").
			Joins("JOIN roles ON users.role = roles.name").
			Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
			Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
			Where("users.id = ? AND permissions.name = ?", claims.UserID, requiredPermission).
			Where("users.deleted_at IS NULL AND roles.deleted_at IS NULL AND permissions.deleted_at IS NULL").
			Count(&count).Error

		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "Database error during authorization")
		}

		if count == 0 {
			return response.Error(c, fiber.StatusForbidden, "Access denied: insufficient permissions")
		}

		return c.Next()
	}
}
