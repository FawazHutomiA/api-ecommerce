package config

import (
	"example/pkg/helper"
	"example/pkg/log"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func InitDB() (*sqlx.DB, error) {
	logger := log.New()

	dbHost := helper.GetENV("DB_HOST")
	dbUser := helper.GetENV("DB_USER")
	dbPassword := helper.GetENV("DB_PASSWORD")
	dbName := helper.GetENV("DB_NAME")
	dbPort := helper.GetENV("DB_PORT")
	dbSSLMode := helper.GetENV("DB_SSLMODE")
	dbTimezone := helper.GetENV("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode, dbTimezone)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		logger.Infof("Failed to open database: %v", err)
		return nil, err
	}

	// Check if the connection to the database is alive
	if err = db.Ping(); err != nil {
		logger.Infof("Failed to ping database: %v", err)
		return nil, err
	}

	logger.Infof("Successfully connected to the database")
	return db, nil
}
