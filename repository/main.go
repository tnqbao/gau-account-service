package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

var repository *Repository

func InitRepository(db *gorm.DB) *Repository {
	repository = &Repository{
		db: db,
	}
	if repository.db == nil {
		panic("database connection is nil")
	}
	return repository
}

func GetRepository() *Repository {
	if repository == nil {
		panic("repository not initialized")
	}
	return repository
}
