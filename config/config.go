package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/marktrs/simple-todo/logger"
)

// Config func to get env value
func Config(key string) string {
	log := logger.Log

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Warn().AnErr("error", err).Msg("Error loading .env file, using system env")
	}

	return os.Getenv(key)
}
