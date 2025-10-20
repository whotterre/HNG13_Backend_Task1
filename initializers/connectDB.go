package initializers

import (
	"fmt"
	"log"
	"os"
	"task_one/config"
	"task_one/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(conf *config.Config) (*gorm.DB, error) {
	dsn := conf.DBUrl
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		return nil, fmt.Errorf("no database DSN provided; set DB_URL or DATABASE_URL")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}

	log.Println("Successfully connected to PostgresDB (via GORM)")
	return db, nil
}

func DoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.StringEntry{})
	if err != nil {
		log.Println("Failed to perform migrations")
		return err
	}
	return nil
}
