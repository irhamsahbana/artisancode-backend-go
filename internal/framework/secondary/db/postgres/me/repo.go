package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type meRepo struct {
	db *sqlx.DB
}

var _ portsRepo.MeRepository = &meRepo{}

type Config struct {
	DB *sqlx.DB
}

func NewMeRepository(cfg Config) portsRepo.MeRepository {
	return &meRepo{db: cfg.DB}
}
