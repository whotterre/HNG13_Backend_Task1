package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"task_one/dto"
)

type NaturalLanguageParser interface {
	ParseQuery(query string) (*dto.FilterByCriteriaData, *dto.InterpretedQuery, error)
}

type naturalLanguageParser struct{}

func NewNaturalLanguageParser() NaturalLanguageParser {
	return &naturalLanguageParser{}
}

func (p *naturalLanguageParser) ParseQuery(query string) (*dto.FilterByCriteriaData, *dto.InterpretedQuery, error) {
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	filters := &dto.FilterByCriteriaData{}
	parsedFilters := make(map[string]any)

	// Parse different patterns
	err := p.parsePatterns(normalizedQuery, filters, parsedFilters)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse natural language query: %v", err)
	}

	// Validate for conflicts
	if err := p.validateFilters(filters); err != nil {
		return nil, nil, fmt.Errorf("query parsed but resulted in conflicting filters: %v", err)
	}

	interpretedQuery := &dto.InterpretedQuery{
		Original:      query,
		ParsedFilters: parsedFilters,
	}

	return filters, interpretedQuery, nil
}

func (p *naturalLanguageParser) parsePatterns(query string, filters *dto.FilterByCriteriaData, parsedFilters map[string]any) error {
	words := strings.Fields(query)

	// Pattern 1: "single word" -> word_count = 1
	if p.containsPhrase(words, "single", "word") {
		wordCount := 1
		filters.WordCount = &wordCount
		parsedFilters["word_count"] = wordCount
	}

	// Pattern 2: "palindromic" -> is_palindrome = true
	if p.containsWord(words, "palindromic") {
		isPalindrome := true
		filters.IsPalindrome = &isPalindrome
		parsedFilters["is_palindrome"] = isPalindrome
	}

	// Pattern 3: "longer than X characters" -> min_length = X + 1
	if longerThanPattern := regexp.MustCompile(`longer than (\d+) characters?`); longerThanPattern.MatchString(query) {
		matches := longerThanPattern.FindStringSubmatch(query)
		if len(matches) > 1 {
			if length, err := strconv.Atoi(matches[1]); err == nil {
				minLength := length + 1
				filters.MinLength = &minLength
				parsedFilters["min_length"] = minLength
			}
		}
	}

	// Pattern 4: "containing the letter X" or "contain the letter X" -> contains_character = X
	if containLetterPattern := regexp.MustCompile(`contain(?:ing)? the letter ([a-z])`); containLetterPattern.MatchString(query) {
		matches := containLetterPattern.FindStringSubmatch(query)
		if len(matches) > 1 {
			character := matches[1]
			filters.ContainsCharacter = &character
			parsedFilters["contains_character"] = character
		}
	}

	// Pattern 5: "containing the first vowel" -> contains_character = a
	if p.containsPhrase(words, "first", "vowel") && p.containsWord(words, "containing") {
		character := "a"
		filters.ContainsCharacter = &character
		parsedFilters["contains_character"] = character
	}

	return nil
}

func (p *naturalLanguageParser) containsWord(words []string, target string) bool {
	for _, word := range words {
		if word == target {
			return true
		}
	}
	return false
}

func (p *naturalLanguageParser) containsPhrase(words []string, phrase ...string) bool {
	if len(phrase) == 0 {
		return false
	}

	for i := 0; i <= len(words)-len(phrase); i++ {
		match := true
		for j, word := range phrase {
			if words[i+j] != word {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func (p *naturalLanguageParser) validateFilters(filters *dto.FilterByCriteriaData) error {
	// Check for min_length > max_length
	if filters.MinLength != nil && filters.MaxLength != nil {
		if *filters.MinLength > *filters.MaxLength {
			return fmt.Errorf("min_length (%d) cannot be greater than max_length (%d)", *filters.MinLength, *filters.MaxLength)
		}
	}

	return nil
}
