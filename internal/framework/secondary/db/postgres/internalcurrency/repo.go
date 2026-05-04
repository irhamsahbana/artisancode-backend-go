package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type internalCurrencyRepo struct {
	db *sqlx.DB
}

type Config struct {
	DB *sqlx.DB
}

var _ portsRepo.InternalCurrencyRepository = &internalCurrencyRepo{}

func NewInternalCurrencyRepository(cfg Config) portsRepo.InternalCurrencyRepository {
	return &internalCurrencyRepo{db: cfg.DB}
}
