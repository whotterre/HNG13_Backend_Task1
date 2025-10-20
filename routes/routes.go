package routes

import (
	"database/sql"
	"task_one/handlers"
	"task_one/repository"
	"task_one/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	stringRepo := repository.NewStringRepository(db)
	stringService := services.NewStringService(stringRepo)
 	stringHandler := handlers.NewStringsHandler(*stringService)
	router.POST("/strings", stringHandler.CreateNewString)
}