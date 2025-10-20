package services

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// getLength returns the number of characters in the string
func getLength(value string) int {
	return len(value)
}

// getIsPalindrome checks if the string reads the same forwards and backwards (case-insensitive, ignoring spaces)
func getIsPalindrome(value string) bool {
	lower := strings.ToLower(value)
	clean := strings.ReplaceAll(lower, " ", "")
	return clean == reverseString(clean)
}

// reverseString reverses the input string
func reverseString(value string) string {
	var b strings.Builder
	for i := len(value) - 1; i >= 0; i-- {
		b.WriteByte(value[i])
	}
	return b.String()
}

// getUniqueCharsCount returns the number of distinct characters in the string
func getUniqueCharsCount(value string) int {
	seen := make(map[rune]struct{})
	for _, r := range value {
		seen[r] = struct{}{}
	}
	return len(seen)
}

// getHash computes the SHA-256 hash of the string
func getHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", sum[:])
}

// getCharFreqMap returns a map of character to occurrence count
func getCharFreqMap(value string) map[string]int {
	freqMap := make(map[string]int)
	normalized := strings.TrimSpace(strings.ToLower(value))
	for _, r := range normalized {
		ch := string(r)
		freqMap[ch]++
	}
	return freqMap
}

// getWordCount returns the number of words in the string
func getWordCount(value string) int {
	words := strings.Fields(value)
	return len(words)
}
