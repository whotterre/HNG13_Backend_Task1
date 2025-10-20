package services

import (
	// "strings"
	"task_one/dto"
	"task_one/repository"
)

type StringService struct {
	stringRepo repository.StringRepository
}

func NewStringService(stringRepo repository.StringRepository) *StringService {
	return &StringService{
		stringRepo:stringRepo,
	}
}

func (s *StringService) CreateNewString(input dto.CreateNewStringEntryRequest) (dto.CreateNewStringResponse, error){ 
	
	return dto.CreateNewStringResponse{}, nil
}
