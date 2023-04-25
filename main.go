package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/marktrs/simple-todo/database"
	"github.com/marktrs/simple-todo/logger"
	"github.com/marktrs/simple-todo/repository"
	"github.com/marktrs/simple-todo/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log := logger.Log
	app := fiber.New(fiber.Config{
		AppName:               "simple-todo-api",
		DisableStartupMessage: true,
		EnablePrintRoutes:     false,
	})

	// Connect to the database
	database.ConnectDB()

	// Setup routes
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()

	router.SetupRoutes(app, userRepo, taskRepo)

	// Listen from a different goroutine
	go func() {
		log.Info().Msg("server is listening on port 3000...")
		if err := app.Listen(":3000"); err != nil {
			log.Panic().AnErr("error", err).Msg("server failed to start")
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	// This blocks the main thread until an interrupt is received
	<-c
	log.Info().Msg("gracefully shutting down...")

	if err := app.Shutdown(); err != nil {
		log.Panic().AnErr("error", err).Msg("server failed to shutdown gracefully")
	}

	cleanup()
}

func cleanup() {
	log := logger.Log
	log.Info().Msg("running cleanup tasks...")

	// Close the database connection
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Panic().AnErr("error", err).Msg("failed to get sqlDB on closing database connection")
	}

	if err := sqlDB.Close(); err != nil {
		log.Panic().AnErr("error", err).Msg("failed to close database connection")
	}
	log.Info().Msg("server was successful shutdown.")
}
