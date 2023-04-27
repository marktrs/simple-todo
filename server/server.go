package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/marktrs/simple-todo/database"
	"github.com/marktrs/simple-todo/logger"
	"github.com/marktrs/simple-todo/middleware"
	"github.com/rs/zerolog"
)

type Server interface {
	App() *fiber.App
	Start()
}

type server struct {
	app    *fiber.App
	logger *zerolog.Logger
}

func New() Server {
	app := fiber.New(fiber.Config{
		AppName:               "simple-todo-api",
		DisableStartupMessage: true,
		EnablePrintRoutes:     false,
		ErrorHandler:          middleware.HandleHTTPError,
	})

	// Enable HTTP caching to intercept responses and cache them
	// app.Use(cache.New())

	// Setup recover middleware to recovers from panics anywhere
	// in the stack chain and handles the control to the centralized ErrorHandler.
	app.Use(recover.New())

	// Setup CORS middleware configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Logger middleware for all routes
	app.Use(middleware.HandleFiberCtxLogger)

	return &server{
		app:    app,
		logger: &logger.Log,
	}
}

func (s *server) App() *fiber.App {
	return s.app
}

func (s *server) Start() {
	go func() {
		s.logger.Info().Msg("server is listening on port 3000...")
		if err := s.app.Listen(":3000"); err != nil {
			s.logger.Panic().AnErr("error", err).Msg("server failed to start")
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	// This blocks the main thread until an interrupt is received
	<-c
	s.logger.Info().Msg("gracefully shutting down...")

	if err := s.app.Shutdown(); err != nil {
		s.logger.Panic().AnErr("error", err).Msg("server failed to shutdown gracefully")
	}

	s.cleanup()
}

func (s *server) cleanup() {
	s.logger.Info().Msg("running cleanup tasks...")

	// Close the database connection
	sqlDB, err := database.DB.DB()
	if err != nil {
		s.logger.Panic().AnErr("error", err).Msg("failed to get sqlDB on closing database connection")
	}

	if err := sqlDB.Close(); err != nil {
		s.logger.Panic().AnErr("error", err).Msg("failed to close database connection")
	}

	s.logger.Info().Msg("server was successful shutdown.")
}
