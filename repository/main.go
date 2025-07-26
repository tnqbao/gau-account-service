package repository

import (
	"github.com/tnqbao/gau-account-service/infra"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
	//cacheDb *redis.Client
}

var repository *Repository

func InitRepository(infra *infra.Infra) *Repository {
	repository = &Repository{
		db: infra.Postgres.DB,
		//cacheDb: infra.Redis.Client,
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
