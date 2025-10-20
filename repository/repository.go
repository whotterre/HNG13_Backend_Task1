package repository

import "database/sql"

type StringRepository interface {

}

type stringRepository struct {
	db *sql.DB
}

func NewStringRepository(db *sql.DB) StringRepository {
	return &stringRepository{db: db}
}





