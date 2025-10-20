package repository

import (
	"gorm.io/gorm"
)

type StringRepository interface {
}

type stringRepository struct {
	db *gorm.DB
}

func NewStringRepository(db *gorm.DB) StringRepository {
	return &stringRepository{db: db}
}
