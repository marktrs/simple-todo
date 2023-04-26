package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/marktrs/simple-todo/handler"
	"github.com/marktrs/simple-todo/middleware"
	"github.com/marktrs/simple-todo/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

// SetupRoutes setup router api
func SetupRoutes(
	app *fiber.App,
	userRepo repository.UserRepository,
	taskRepo repository.TaskRepository,
) {
	validator := validator.New()

	// Metrics
	app.Get("/metrics", monitor.New(monitor.Config{Title: "Simple-TODO API Metrics"}))

	api := app.Group("/api")
	// Health check
	api.Get("/health", handler.HealthCheck)

	// Auth
	authHandler := handler.NewAuthHandler(validator, userRepo)
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)

	// User
	userHandler := handler.NewUserHandler(validator, userRepo)
	user := api.Group("/users")
	user.Post("/", userHandler.CreateUser)

	// Task
	taskHandler := handler.NewTaskHandler(validator, taskRepo)
	task := api.Group("/tasks")
	task.Use(middleware.Protected())
	task.Get("/", taskHandler.GetAllTasks)
	task.Post("/", taskHandler.CreateTask)
	task.Put("/:id", taskHandler.UpdateTask)
	task.Delete("/:id", taskHandler.DeleteTask)
}
