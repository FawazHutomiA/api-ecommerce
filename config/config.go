package config

import (
	"example/internal/helper"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dbHost := helper.GetENV("DB_HOST")
	dbUser := helper.GetENV("DB_USER")
	dbPassword := helper.GetENV("DB_PASSWORD")
	dbName := helper.GetENV("DB_NAME")
	dbPort := helper.GetENV("DB_PORT")
	dbSSLMode := helper.GetENV("DB_SSLMODE")
	dbTimezone := helper.GetENV("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode, dbTimezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
