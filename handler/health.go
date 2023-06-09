package handler

import (
	"github.com/gofiber/fiber/v2"
)

// HealthCheck handle api status
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "API is running"})
}
