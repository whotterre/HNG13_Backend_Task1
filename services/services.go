package services

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

func getLength(value string) int {
	return len(value)
}

func getIsPalindrome(value string) bool {
	lower := strings.ToLower(value)
	clean := strings.ReplaceAll(lower, " ", "")
	return clean == reverseString(clean)
}

func reverseString(value string) string {
	var b strings.Builder
	for i := len(value) - 1; i >= 0; i-- {
		b.WriteByte(value[i])
	}
	return b.String()
}

func getUniqueCharsCount(value string) int {
	seen := make(map[rune]struct{})
	for _, r := range value {
		seen[r] = struct{}{}
	}
	return len(seen)
}

// Get SHA256 hash for the string
func getHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", sum[:])
}

func getCharFreqMap(value string) map[string]int {
	freqMap := make(map[string]int)
	normalized := strings.TrimSpace(strings.ToLower(value))
	for _, r := range normalized {
		ch := string(r)
		freqMap[ch]++
	}
	return freqMap
}

func getWordCount(value string) int {
	words := strings.Fields(value)
	return len(words)
}
