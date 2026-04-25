package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type internalClientRepo struct {
	db *sqlx.DB
}

type InternalClientRepositoryConfig struct {
	DB *sqlx.DB
}

var _ portsRepo.InternalClientRepository = &internalClientRepo{}

func NewInternalClientRepository(cfg InternalClientRepositoryConfig) portsRepo.InternalClientRepository {
	return &internalClientRepo{db: cfg.DB}
}
