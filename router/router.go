package router

import (
	"github.com/marktrs/simple-todo/handler"
	"github.com/marktrs/simple-todo/logger"
	"github.com/marktrs/simple-todo/middleware"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	app.Use(recover.New())
	app.Use(cors.New())

	// Logger middleware for all routes
	app.Use(adaptor.HTTPMiddleware(logger.HttpLogger))

	// Metrics
	app.Get("/metrics", monitor.New(monitor.Config{Title: "Simple-TODO API Metrics"}))

	api := app.Group("/api")
	// Health check
	api.Get("/health", handler.HealthCheck)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)

	// User
	user := api.Group("/users")
	user.Post("/", handler.CreateUser)

	// Task
	task := api.Group("/tasks")
	task.Use(middleware.Protected())
	task.Get("/", handler.GetAllTasks)
	task.Post("/", handler.CreateTask)
	task.Put("/:id", handler.UpdateTask)
	task.Delete("/:id", handler.DeleteTask)
}
