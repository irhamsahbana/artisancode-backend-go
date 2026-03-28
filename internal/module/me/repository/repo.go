package repository

import (
	portsRepo "codebase-app/internal/ports/repository"

	"github.com/jmoiron/sqlx"
)

type meRepo struct {
	db *sqlx.DB
}

var _ portsRepo.MeRepository = &meRepo{}

type MeRepositoryConfig struct {
	DB *sqlx.DB
}

func NewMeRepository(cfg MeRepositoryConfig) portsRepo.MeRepository {
	return &meRepo{db: cfg.DB}
}
