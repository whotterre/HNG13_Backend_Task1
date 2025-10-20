package handlers

import (
	"net/http"
	"task_one/dto"
	"task_one/services"

	"github.com/gin-gonic/gin"
)

type StringsHandler struct {
	stringsService services.StringService
}

func NewStringsHandler(stringService services.StringService) *StringsHandler {
	return &StringsHandler{
		stringsService: stringService,
	}
}



func (h *StringsHandler) CreateNewString(c *gin.Context){
	// Bind body 
	var req dto.CreateNewStringEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read request body",
		})
	}

	if len(req.Value) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing value in request",
		})
	}

	response, err := h.stringsService.CreateNewString(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":"failed to get response from service layer",
		})
	}


	c.JSON(http.StatusCreated, response)
}