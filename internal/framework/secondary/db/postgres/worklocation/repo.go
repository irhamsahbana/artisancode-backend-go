package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type workLocationRepo struct {
	db *sqlx.DB
}

var _ portsRepo.WorkLocationRepository = &workLocationRepo{}

type Config struct {
	DB *sqlx.DB
}

func NewWorkLocationRepository(cfg Config) portsRepo.WorkLocationRepository {
	return &workLocationRepo{db: cfg.DB}
}
