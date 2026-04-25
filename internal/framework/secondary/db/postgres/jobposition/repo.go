package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type jobPositionRepo struct {
	db *sqlx.DB
}

var _ portsRepo.JobPositionRepository = &jobPositionRepo{}

type JobPositionRepositoryConfig struct {
	DB *sqlx.DB
}

func NewJobPositionRepository(cfg JobPositionRepositoryConfig) portsRepo.JobPositionRepository {
	return &jobPositionRepo{db: cfg.DB}
}
