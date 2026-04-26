package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type internalQuotationRepo struct {
	db *sqlx.DB
}

type Config struct {
	DB *sqlx.DB
}

var _ portsRepo.InternalQuotationRepository = &internalQuotationRepo{}

func NewInternalQuotationRepository(cfg Config) portsRepo.InternalQuotationRepository {
	return &internalQuotationRepo{db: cfg.DB}
}
