package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
	Port  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		// Try parent directory
		err = godotenv.Load("../.env")
		if err != nil {
			log.Println("No .env file found, using system environment variables")
		}
	}

	dbUrl := getEnv("DATABASE_URL", "")
	if dbUrl == "" {
		dbUrl = getEnv("DB_URL", "")
	}

	return &Config{
		DBUrl: dbUrl,
		Port:  getEnv("PORT", "4000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
