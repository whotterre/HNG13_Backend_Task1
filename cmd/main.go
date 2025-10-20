package main

import (
	"fmt"
	"log"
	"task_one/config"
	"task_one/initializers"
	"task_one/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	router := gin.Default()

	db, err := initializers.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL database: %v", err)
	}

	routes.SetupRoutes(router, db)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start http server: %v", err)
	}
}
