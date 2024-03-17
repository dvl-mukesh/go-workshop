package database

import (
	"fmt"
	"log"

	"github.com/dvl-mukesh/go-workshop/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(env *config.Environment) (*gorm.DB, error) {
	log.Println("Setting up new db connection")

	connectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", env.DbHost, env.DbPort, env.DbUserName, env.DbPassword, env.DbName)

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
