package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/marktrs/simple-todo/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:               "simple-todo-api",
		DisableStartupMessage: true,
		EnablePrintRoutes:     false,
	})

	logFileDir := strings.Join([]string{"./temp/", time.Now().Format("2006-01-02_15:04:05"), ".log"}, "")

	file, err := os.OpenFile(logFileDir, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Output:     file,
		TimeFormat: time.RFC3339Nano,
		Format:     `{"timestamp":"${time}", "request_id":"${locals:requestid}", "status":${status},"latency":"${latency}", "path":"${path}",​ "body": ${body}},`,
	}))
	app.Use(recover.New())
	app.Use(cors.New())

	router.SetupRoutes(app)

	// Listen from a different goroutine
	go func() {
		log.Print("server is listning on port 3000...")
		if err := app.Listen(":3000"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	// This blocks the main thread until an interrupt is received
	_ = <-c
	log.Print("gracefully shutting down...")

	if err := app.Shutdown(); err != nil {
		log.Panic(err)
	}

	log.Print("running cleanup tasks...")

	// TODO: add cleanup tasks

	log.Print("server was successful shutdown.")
}