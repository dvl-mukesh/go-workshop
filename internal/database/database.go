package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {
	fmt.Println("Setting up new db connection")

	dbUserName := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbTable := os.Getenv("DB_TABLE")
	dbPort := os.Getenv("DB_PORT")

	connectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUserName, dbPassword, dbTable)

	db, err := gorm.Open(postgres.Open(connectString))

	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()

	if err != nil {
		return nil, err
	}

	if err := sqlDb.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
