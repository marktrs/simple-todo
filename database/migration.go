package database

import (
	"github.com/marktrs/simple-todo/model"
)

func migrateTables() error {
	// Migrate the schema
	if err := DB.AutoMigrate(
		&model.User{},
		&model.Task{},
	); err != nil {
		return err
	}

	return nil
}
