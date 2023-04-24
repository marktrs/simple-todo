package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/marktrs/simple-todo/config"
	"github.com/marktrs/simple-todo/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB initialize database connection
func ConnectDB() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		log.Fatal("failed to parse database port")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))

	if DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Error),
	}); err != nil {
		log.Fatal("failed to connect database")
	}

	if err := DB.AutoMigrate(&model.Task{}, &model.User{}); err != nil {
		log.Fatal("failed to migrate database")
	}
}
