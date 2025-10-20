package dto

type CreateNewStringEntryRequest struct {
	Value string `json:"value"`
}

type StringProperties struct {
	Length       int          `json:"length"`
	IsPalindrome bool         `json:"is_palindrome"`
	UniqueChars  int          `json:"unique_characters"`
	WordCount    int          `json:"word_count"`
	FreqMap      map[rune]int `json:"character_frequency_map"`
}

type CreateNewStringResponse struct {
	Id         string           `json:"id"`
	Value      int              `json:"value"`
	Properties StringProperties `json:"properties"`
}
