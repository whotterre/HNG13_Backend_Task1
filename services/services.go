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

type StringService struct {
	stringRepo repository.StringRepository
}

func NewStringService(stringRepo repository.StringRepository) *StringService {
	return &StringService{
		stringRepo: stringRepo,
	}
}

func (s *StringService) CreateNewString(input dto.CreateNewStringEntryRequest) (*dto.CreateNewStringResponse, error) {
	// Check for duplicates by value
	if existing, err := s.stringRepo.GetStringByValue(input.Value); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, fmt.Errorf("conflict: string already exists")
	}
	// Compute string details
	stringDetails := models.StringDetails{
		Hash:         getHash(input.Value),
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
			Hash:         stringDetails.Hash,
			FreqMap:      stringDetails.FreqMap,
			CreatedAt:    now,
		},
	}

	return &finalResponse, nil
}
