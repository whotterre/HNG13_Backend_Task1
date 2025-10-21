package initializers

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
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
	// First attempt to connect using the provided DSN (use the DB the user specified)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Check if the error is because the database doesn't exist
		if strings.Contains(err.Error(), "does not exist") {
			log.Printf("Database does not exist, attempting to create it: %v", err)
			// Try to create the database named in the DSN
			if createErr := createDatabase(dsn); createErr != nil {
				log.Printf("Failed to create database: %v", createErr)
				return nil, fmt.Errorf("failed to connect to database and failed to create it: %v", err)
			}

			// Try to connect again after creating the database
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Printf("Failed to connect to database after creation: %v", err)
				return nil, err
			}
		} else {
			log.Printf("Unable to connect to database: %v\n", err)
			return nil, err
		}
	}

	log.Println("Successfully connected to PostgresDB (via GORM)")
	return db, nil
}

func createDatabase(dsn string) error {
	// Parse the DSN to create admin DSN
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %v", err)
	}

	// Create admin DSN (connect to 'postgres' database instead)
	adminURL := *parsedURL
	adminURL.Path = "/postgres"
	adminDSN := adminURL.String()

	// Connect to admin database
	adminDB, err := gorm.Open(postgres.Open(adminDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to admin database: %v", err)
	}

	// Determine the database name from the DSN path (use the last path segment)
	dbPath := parsedURL.Path
	// Trim any leading slash
	dbPath = strings.TrimPrefix(dbPath, "/")
	dbName := path.Base(dbPath)
	if dbName == "" {
		return fmt.Errorf("failed to determine database name from DSN path: %s", parsedURL.Path)
	}

	// Quote the database name to safely handle unusual names
	createQuery := fmt.Sprintf("CREATE DATABASE %q", dbName)
	if err := adminDB.Exec(createQuery).Error; err != nil {
		// Check if database already exists (race condition)
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("failed to create database: %v", err)
		}
	}

	log.Printf("Successfully created database: %s", dbName)
	return nil
}


func DoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.StringEntry{})
	if err != nil {
		log.Println("Failed to perform migrations")
		return err
	}
	return nil
}
