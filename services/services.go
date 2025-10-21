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

func (s *stringService) FilterByCriteria(input dto.FilterByCriteriaData) (*dto.FilterByCriteriaResponse, error){
	// Filter by data
	stringData, err := s.stringRepo.FilterByCriteria(input)
	if err != nil {
		return nil, err		
	}

	// Transform []models.StringEntry to []dto.GetStringByValueResponse
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

	// Convert input struct to map[string]any
	filtersMap := make(map[string]any)
	inputJSON, _ := json.Marshal(input)
	json.Unmarshal(inputJSON, &filtersMap)
	
	response := dto.FilterByCriteriaResponse{
		Data: transformedData,
		Count: len(transformedData),
		FiltersApplied: filtersMap,
	}
	return &response, nil
}