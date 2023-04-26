package main

import (
	"github.com/marktrs/simple-todo/database"
	"github.com/marktrs/simple-todo/repository"
	"github.com/marktrs/simple-todo/router"
	"github.com/marktrs/simple-todo/server"
)

func main() {
	// Create a new fiber app
	srv := server.New()

	// Connect to the database
	database.ConnectDB()

	// Initialize repositories used by the routes
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()

	// Setup routes
	router.SetupRoutes(srv.App(), userRepo, taskRepo)

	// Start the server
	srv.Start()
}
