package database

import (
	"github.com/marktrs/simple-todo/logger"
	"github.com/marktrs/simple-todo/model"
)

func migrateTables() {
	log := logger.Log
	// Migrate the schema
	if err := DB.AutoMigrate(
		&model.User{},
		&model.Task{},
	); err != nil {
		log.Fatal().AnErr("error", err).Msg("failed to migrate database")
	}
}
