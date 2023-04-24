package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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

	host := config.Config("DB_HOST")
	if os.Getenv("DB_HOST") != "" {
		host = os.Getenv("DB_HOST")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		host,
		port,
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
	)

	if DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Error),
	}); err != nil {
		log.Fatal("failed to connect database")
	}

	migrateTables()
	setConnectionPool()
}

func migrateTables() {
	// Migrate the schema
	if err := DB.AutoMigrate(&model.User{}, &model.Task{}); err != nil {
		log.Fatal("failed to migrate database")
	}
}

func setConnectionPool() {
	// Get the underlying sql.DB object of the gorm.DB object to use its functions
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("failed to get database connection")
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
}
