package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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

func (h *StringsHandler) CreateNewString(c *gin.Context) {
	// Bind body
	var req dto.CreateNewStringEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If the field "value" was the wrong type, return 422
		if ute, ok := err.(*json.UnmarshalTypeError); ok && strings.EqualFold(ute.Field, "value") {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Invalid data type for \"value\"; must be string"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to read request body"})
		return
	}

	if len(req.Value) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing value in request",
		})
		return
	}

	response, err := h.stringsService.CreateNewString(req)
	if err != nil {
		// Map duplicate error to 409 Conflict
		if strings.Contains(err.Error(), "conflict") {
			c.JSON(http.StatusConflict, gin.H{"message": "String already exists in the system"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get response from service layer"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *StringsHandler) GetStringByValue(c *gin.Context) {
	stringValue := c.Param("string_value")

	response, err := h.stringsService.GetStringByValue(stringValue)
	if err != nil {
		// Map not found error to 404
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"message": "String does not exist in the system"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve string"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *StringsHandler) FilterByCriteria(c *gin.Context) {
	isPalindrome := c.Query("is_palindrome")
	minLength := c.Query("min_length")
	maxLength := c.Query("max_length")
	wordCount := c.Query("word_count")
	containsCharacter := c.Query("contains_character")

	// Convert query elements to their actual data types
	var minLengthInt int
	if minLength != "" {
		val, err := strconv.Atoi(minLength)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad query"})
			return
		}
		minLengthInt = val
	}

	var maxLengthInt int
	if maxLength != "" {
		val, err := strconv.Atoi(maxLength)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad query"})
			return
		}
		maxLengthInt = val
	}

	var wordCountInt int
	if wordCount != "" {
		val, err := strconv.Atoi(wordCount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad query"})
			return
		}
		wordCountInt = val
	}

	// Check if isPalindrome is a valid boolean
	isPalindromeBool := false
	if isPalindrome != "" {
		isPalindromeBool = (isPalindrome == "true")
	}

	// Check if containsCharacter is provided
	containsCharacterStr := containsCharacter
	// Populate the dto
	input := dto.FilterByCriteriaData{
		IsPalindrome:      isPalindromeBool,
		MinLength:         minLengthInt,
		MaxLength:         maxLengthInt,
		WordCount:         wordCountInt,
		ContainsCharacter: containsCharacterStr,
	}

	response, err := h.stringsService.FilterByCriteria(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve strings"})
		return
	}

	c.JSON(http.StatusOK, response)
}
