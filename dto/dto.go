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
	SHA256Hash   string         `json:"sha256_hash"`
	FreqMap      map[string]int `json:"character_frequency_map"`
}

type CreateNewStringResponse struct {
	Id         string           `json:"id"`
	Value      string           `json:"value"`
	Properties StringProperties `json:"properties"`
	CreatedAt    time.Time      `json:"created_at"`
}

type GetStringByValueResponse struct {
	Id         string           `json:"id"`
	Value      string           `json:"value"`
	Properties StringProperties `json:"properties"`
	CreatedAt    string      `json:"created_at"`
}