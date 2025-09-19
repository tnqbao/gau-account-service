package repository

import (
	"github.com/tnqbao/gau-account-service/shared/infra"
	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
	//cacheDb *redis.Client
}

var repository *Repository

func InitRepository(infra *infra.Infra) *Repository {
	repository = &Repository{
		Db: infra.Postgres.DB,
		//cacheDb: infra.Redis.Client,
	}
	if repository.Db == nil {
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

// Transaction support methods
func (r *Repository) BeginTransaction() *gorm.DB {
	return r.Db.Begin()
}

func (r *Repository) WithTransaction(tx *gorm.DB) *Repository {
	return &Repository{
		Db: tx,
	}
}
