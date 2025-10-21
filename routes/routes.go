package routes

import (
	"task_one/handlers"
	"task_one/repository"
	"task_one/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	stringRepo := repository.NewStringRepository(db)
	stringService := services.NewStringService(stringRepo)
	stringHandler := handlers.NewStringsHandler(*stringService)
	// Routes
	router.POST("/strings", stringHandler.CreateNewString)
	router.GET("/strings/:string_value", stringHandler.GetStringByValue)
}
