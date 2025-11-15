package middleware

import (
	"strings"

	"mastercard-backend/internal/config"

	"github.com/gofiber/fiber/v2"
)

// CORSMiddleware handles CORS headers
func CORSMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		origins := strings.Split(config.AppConfig.CORSAllowedOrigins, ",")
		origin := c.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range origins {
			if strings.TrimSpace(allowedOrigin) == origin || strings.TrimSpace(allowedOrigin) == "*" {
				allowed = true
				break
			}
		}

		if allowed {
			c.Set("Access-Control-Allow-Origin", origin)
		}

		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Methods", config.AppConfig.CORSAllowedMethods)
		c.Set("Access-Control-Allow-Headers", config.AppConfig.CORSAllowedHeaders)

		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}

