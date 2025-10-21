package repository

import (
	"errors"
	"task_one/dto"
	"task_one/models"

	"gorm.io/gorm"
)

type StringRepository interface {
	CreateNewStringRecord(stringData models.StringEntry) (*models.StringEntry, error)
	GetStringByValue(value string) (*models.StringEntry, error)
	GetStringById(id string) (*models.StringEntry, error)
	FilterByCriteria(input dto.FilterByCriteriaData) (*[]models.StringEntry, error)
	DeleteStringValue(hash string) error
}

type stringRepository struct {
	db *gorm.DB
}

func NewStringRepository(db *gorm.DB) StringRepository {
	return &stringRepository{db: db}
}

func (r stringRepository) GetStringByValue(value string) (*models.StringEntry, error) {
	var entry models.StringEntry
	err := r.db.Where("value = ?", value).First(&entry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r stringRepository) GetStringById(id string) (*models.StringEntry, error) {
	var entry models.StringEntry
	err := r.db.Where("id = ?", id).First(&entry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r stringRepository) CreateNewStringRecord(stringData models.StringEntry) (*models.StringEntry, error) {
	if err := r.db.Create(&stringData).Error; err != nil {
		return nil, err
	}
	return &stringData, nil
}

func (r stringRepository) FilterByCriteria(input dto.FilterByCriteriaData) (*[]models.StringEntry, error) {
	var entries []models.StringEntry
	query := r.db.Model(&models.StringEntry{})

	// Add conditions only if the filter values are provided
	if input.IsPalindrome != nil {
		query = query.Where("is_palindrome = ?", *input.IsPalindrome)
	}
	if input.MinLength != nil {
		query = query.Where("length >= ?", *input.MinLength)
	}
	if input.MaxLength != nil {
		query = query.Where("length <= ?", *input.MaxLength)
	}
	if input.WordCount != nil {
		query = query.Where("word_count = ?", *input.WordCount)
	}

	// Check if JSONB field contains a specific key
	if input.ContainsCharacter != nil {
		query = query.Where("character_frequency_map -> ? IS NOT NULL", *input.ContainsCharacter)
	}

	if err := query.Find(&entries).Error; err != nil {
		return nil, err
	}
	return &entries, nil
}

func (r stringRepository) DeleteStringValue(hash string) error {
	if err := r.db.Where("id = ?", hash).Delete(&models.StringEntry{}).Error; err != nil {
		return err
	}
	return nil
}
