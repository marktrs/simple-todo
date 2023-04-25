package database

import (
	"database/sql"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/marktrs/simple-todo/config"
	"github.com/marktrs/simple-todo/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	dbLogger "gorm.io/gorm/logger"
)

// ConnectDB initialize database connection
func ConnectDB() {
	var err error
	log := logger.Log

	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("failed to parse database port")
	}

	host := config.Config("DB_HOST")
	if os.Getenv("DB_HOST") != "" {
		host = os.Getenv("DB_HOST")
	}

	dsn := strings.Join([]string{
		"host=", host,
		" port=", strconv.Itoa(int(port)),
		" user=", config.Config("DB_USER"),
		" password=", config.Config("DB_PASSWORD"),
		" dbname=", config.Config("DB_NAME"),
	}, "")

	if DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         dbLogger.Default.LogMode(dbLogger.Silent), // Disable DB Logger, only show error message on system logger
	}); err != nil {
		log.Fatal().AnErr("error", err).Msg("failed to connect database")
		return
	}

	if err := setConnectionPool(); err != nil {
		log.Fatal().AnErr("error", err).Msg("failed to set connection pool")
		return
	}

	if err := migrateTables(); err != nil {
		log.Fatal().AnErr("error", err).Msg("failed to migrate database")
		return
	}

	if err := seedUserData(); err != nil {
		log.Fatal().AnErr("error", err).Msg("failed to seed user data")
		return
	}
}

func ConnectExistingSQL(sqlDB *sql.DB) {
	var err error
	log := logger.Log
	if DB, err = gorm.Open(postgres.New(postgres.Config{Conn: sqlDB})); err != nil {
		log.Fatal().AnErr("error", err).Msg("failed to connect database")
		return
	}
}

func setConnectionPool() error {
	log := logger.Log
	// Get the underlying sql.DB object of the gorm.DB object to use its functions
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("failed to get database connection")
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}
