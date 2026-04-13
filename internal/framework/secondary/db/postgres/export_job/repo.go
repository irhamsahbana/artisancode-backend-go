package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type exportJobRepo struct {
	db *sqlx.DB
}

var _ portsRepo.ExportJobRepository = &exportJobRepo{}

type ExportJobRepositoryConfig struct {
	DB *sqlx.DB
}

func NewExportJobRepository(cfg ExportJobRepositoryConfig) portsRepo.ExportJobRepository {
	return &exportJobRepo{db: cfg.DB}
}
