package main

import (
	"log"
	"net/http"
	"task_one/initializers"
	"task_one/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	conn, err := initializers.ConnectDB()
	if err != nil {
		log.Printf("Failed to connect to SQLite3 db")
		return 
	}

	routes.SetupRoutes(router, conn)

	addr := ":4000"
	
	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("Failed to start http server because", err)
	}

	router.Run()
}