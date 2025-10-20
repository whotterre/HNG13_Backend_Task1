package dto

import "time"

type CreateNewStringEntryRequest struct {
	Value string `json:"value"`
}

type StringProperties struct {
	Length       int            `json:"length"`
	IsPalindrome bool           `json:"is_palindrome"`
	UniqueChars  int            `json:"unique_characters"`
	WordCount    int            `json:"word_count"`
	FreqMap      map[string]int `json:"character_frequency_map"`
	Hash         string         `json:"sha256_hash"`
	CreatedAt    time.Time      `json:"created_at"`
}

type CreateNewStringResponse struct {
	Id         string           `json:"id"`
	Value      string           `json:"value"`
	Properties StringProperties `json:"properties"`
}
