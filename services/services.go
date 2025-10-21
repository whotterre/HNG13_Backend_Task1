package services

import (
	"encoding/json"
	"fmt"
	"log"
	"task_one/dto"
	"task_one/models"
	"task_one/repository"
	"time"
)

type StringService interface {
	CreateNewString(input dto.CreateNewStringEntryRequest) (*dto.CreateNewStringResponse, error)
	GetStringByValue(value string) (*dto.GetStringByValueResponse, error)
	FilterByCriteria(input dto.FilterByCriteriaData) (*dto.FilterByCriteriaResponse, error)
	FilterByNaturalLanguage(input dto.FilterByNaturalLanguageRequest) (*dto.FilterByNaturalLanguageResponse, error)
	DeleteStringEntry(value string) error
}

type stringService struct {
	stringRepo repository.StringRepository
}

func NewStringService(stringRepo repository.StringRepository) StringService {
	return &stringService{
		stringRepo: stringRepo,
	}
}

func (s *stringService) CreateNewString(input dto.CreateNewStringEntryRequest) (*dto.CreateNewStringResponse, error) {
	// Check for duplicates by value
	if existing, err := s.stringRepo.GetStringByValue(input.Value); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, fmt.Errorf("conflict: string already exists")
	}
	// Compute string details
	stringDetails := models.StringDetails{
		Hash:         GetHash(input.Value),
		Length:       getLength(input.Value),
		IsPalindrome: getIsPalindrome(input.Value),
		UniqueChars:  getUniqueCharsCount(input.Value),
		WordCount:    getWordCount(input.Value),
		FreqMap:      getCharFreqMap(input.Value),
	}

	// Prepare DB entry
	freqMapJSON, _ := json.Marshal(stringDetails.FreqMap)
	now := time.Now().UTC()
	stringEntry := models.StringEntry{
		ID:                    stringDetails.Hash,
		Value:                 input.Value,
		Length:                stringDetails.Length,
		IsPalindrome:          stringDetails.IsPalindrome,
		UniqueCharacters:      stringDetails.UniqueChars,
		WordCount:             stringDetails.WordCount,
		SHA256Hash:            stringDetails.Hash,
		CharacterFrequencyMap: freqMapJSON,
		CreatedAt:             now,
	}
	// Persist
	_, err := s.stringRepo.CreateNewStringRecord(stringEntry)
	if err != nil {
		log.Println("Failed to create new string record", err)
		return nil, err
	}

	finalResponse := dto.CreateNewStringResponse{
		Id:    stringEntry.ID,
		Value: input.Value,
		Properties: dto.StringProperties{
			Length:       stringDetails.Length,
			IsPalindrome: stringDetails.IsPalindrome,
			UniqueChars:  stringDetails.UniqueChars,
			WordCount:    stringDetails.WordCount,
			SHA256Hash:   stringDetails.Hash,
			FreqMap:      stringDetails.FreqMap,
		},
		CreatedAt: now,
	}

	return &finalResponse, nil
}

func (s *stringService) GetStringByValue(value string) (*dto.GetStringByValueResponse, error) {
	// Generate SHA256 sum
	stringHash := GetHash(value)

	stringData, err := s.stringRepo.GetStringById(stringHash)
	if err != nil {
		return nil, err
	}

	// Return not found error if string doesn't exist
	if stringData == nil {
		return nil, fmt.Errorf("not found: string does not exist in the system")
	}

	var freqMap map[string]int
	if err := json.Unmarshal(stringData.CharacterFrequencyMap, &freqMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal frequency map: %v", err)
	}

	response := dto.GetStringByValueResponse{
		Id:    stringData.ID,
		Value: stringData.Value,
		Properties: dto.StringProperties{
			Length:       stringData.Length,
			IsPalindrome: stringData.IsPalindrome,
			UniqueChars:  stringData.UniqueCharacters,
			WordCount:    stringData.WordCount,
			FreqMap:      freqMap,
		},
		CreatedAt: stringData.CreatedAt.Format(time.RFC3339),
	}
	return &response, nil
}

func (s *stringService) FilterByCriteria(input dto.FilterByCriteriaData) (*dto.FilterByCriteriaResponse, error) {
	stringData, err := s.stringRepo.FilterByCriteria(input)
	if err != nil {
		return nil, err
	}

	var transformedData []dto.GetStringByValueResponse
	for _, entry := range *stringData {
		var freqMap map[string]int
		if err := json.Unmarshal(entry.CharacterFrequencyMap, &freqMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal frequency map: %v", err)
		}

		transformedData = append(transformedData, dto.GetStringByValueResponse{
			Id:    entry.ID,
			Value: entry.Value,
			Properties: dto.StringProperties{
				Length:       entry.Length,
				IsPalindrome: entry.IsPalindrome,
				UniqueChars:  entry.UniqueCharacters,
				WordCount:    entry.WordCount,
				FreqMap:      freqMap,
			},
			CreatedAt: entry.CreatedAt.Format(time.RFC3339),
		})
	}

	// Build filters_applied with only non-nil values
	filtersMap := make(map[string]any)
	if input.IsPalindrome != nil {
		filtersMap["is_palindrome"] = *input.IsPalindrome
	}
	if input.MinLength != nil {
		filtersMap["min_length"] = *input.MinLength
	}
	if input.MaxLength != nil {
		filtersMap["max_length"] = *input.MaxLength
	}
	if input.WordCount != nil {
		filtersMap["word_count"] = *input.WordCount
	}
	if input.ContainsCharacter != nil {
		filtersMap["contains_character"] = *input.ContainsCharacter
	}

	response := dto.FilterByCriteriaResponse{
		Data:           transformedData,
		Count:          len(transformedData),
		FiltersApplied: filtersMap,
	}
	return &response, nil
}

func (s *stringService) FilterByNaturalLanguage(input dto.FilterByNaturalLanguageRequest) (*dto.FilterByNaturalLanguageResponse, error) {
	// Parse the natural language query
	parser := NewNaturalLanguageParser()
	filters, interpretedQuery, err := parser.ParseQuery(input.Query)
	if err != nil {
		return nil, err
	}

	// Use the existing FilterByCriteria method
	criteriaResponse, err := s.FilterByCriteria(*filters)
	if err != nil {
		return nil, err
	}

	// Ensure data is not nil
	data := criteriaResponse.Data
	if data == nil {
		data = []dto.GetStringByValueResponse{}
	}

	// Convert to natural language response
	response := &dto.FilterByNaturalLanguageResponse{
		Data:             data,
		Count:            criteriaResponse.Count,
		InterpretedQuery: *interpretedQuery,
	}

	return response, nil
}

func (s *stringService) DeleteStringEntry(value string) error {
	// Compute the hash
	hashValue := GetHash(value)

	// Check if it exists
	existing, err := s.stringRepo.GetStringById(hashValue)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("not found: string does not exist in the system")
	}

	// Delete the string
	return s.stringRepo.DeleteStringValue(hashValue)
}
