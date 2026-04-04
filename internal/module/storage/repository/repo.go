package repository

import (
	portsRepo "codebase-app/internal/ports/repository"

	"github.com/jmoiron/sqlx"
)

type storageRepo struct {
	db *sqlx.DB
}

var _ portsRepo.StorageRepository = &storageRepo{}

type StorageRepositoryConfig struct {
	DB *sqlx.DB
}

func NewStorageRepository(cfg StorageRepositoryConfig) portsRepo.StorageRepository {
	return &storageRepo{db: cfg.DB}
}
