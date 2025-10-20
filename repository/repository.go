package repository

import (
	"errors"
	"task_one/models"

	"gorm.io/gorm"
)

type StringRepository interface {
	CreateNewStringRecord(stringData models.StringEntry) (*models.StringEntry, error)
	GetStringByValue(value string) (*models.StringEntry, error)
}

type stringRepository struct {
	db *gorm.DB
}

func NewStringRepository(db *gorm.DB) StringRepository {
	return &stringRepository{db: db}
}

func (s stringRepository) GetStringByValue(value string) (*models.StringEntry, error) {
	var entry models.StringEntry
	err := s.db.Where("value = ?", value).First(&entry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (s stringRepository) CreateNewStringRecord(stringData models.StringEntry) (*models.StringEntry, error) {
	if err := s.db.Create(&stringData).Error; err != nil {
		return nil, err
	}
	return &stringData, nil
}
