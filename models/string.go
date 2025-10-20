package models

import (
	"time"

	"gorm.io/datatypes"
)

type StringEntry struct {
	ID                    string         `gorm:"primaryKey;type:text" json:"id"`
	Value                 string         `gorm:"type:text;not null" json:"value"`
	Length                int            `gorm:"not null" json:"length"`
	IsPalindrome          bool           `gorm:"not null" json:"is_palindrome"`
	UniqueCharacters      int            `gorm:"not null" json:"unique_characters"`
	WordCount             int            `gorm:"not null" json:"word_count"`
	SHA256Hash            string         `gorm:"type:text;not null" json:"sha256_hash"`
	CharacterFrequencyMap datatypes.JSON `gorm:"type:jsonb;not null" json:"character_frequency_map"`
	CreatedAt             time.Time      `gorm:"autoCreateTime" json:"created_at"`
}

type StringDetails struct {
	Hash         string         `json:"sha256_hash"`
	Length       int            `json:"length"`
	IsPalindrome bool           `json:"is_palindrome"`
	UniqueChars  int            `json:"unique_characters"`
	WordCount    int            `json:"word_count"`
	FreqMap      map[string]int `json:"character_frequency_map"`
}
