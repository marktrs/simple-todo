package router

import (
	"github.com/marktrs/simple-todo/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	app.Get("/metrics", monitor.New(monitor.Config{Title: "Simple-TODO API Metrics"}))

	api := app.Group("/api", logger.New())
	api.Get("/health", handler.HealthCheck)
}
