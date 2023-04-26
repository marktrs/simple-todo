package router

import (
	"github.com/marktrs/simple-todo/handler"
	"github.com/marktrs/simple-todo/logger"
	"github.com/marktrs/simple-todo/middleware"
	"github.com/marktrs/simple-todo/repository"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupRoutes setup router api
func SetupRoutes(
	app *fiber.App,
	userRepo repository.UserRepository,
	taskRepo repository.TaskRepository,
) {
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Logger middleware for all routes
	app.Use(adaptor.HTTPMiddleware(logger.HttpLogger))

	// Metrics
	app.Get("/metrics", monitor.New(monitor.Config{Title: "Simple-TODO API Metrics"}))

	api := app.Group("/api")
	// Health check
	api.Get("/health", handler.HealthCheck)

	// Auth
	authHandler := handler.NewAuthHandler(userRepo)
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)

	// User
	userHandler := handler.NewUserHandler(userRepo)
	user := api.Group("/users")
	user.Post("/", userHandler.CreateUser)

	// Task
	taskHandler := handler.NewTaskHandler(taskRepo)
	task := api.Group("/tasks")
	task.Use(middleware.Protected())
	task.Get("/", taskHandler.GetAllTasks)
	task.Post("/", taskHandler.CreateTask)
	task.Put("/:id", taskHandler.UpdateTask)
	task.Delete("/:id", taskHandler.DeleteTask)
}
