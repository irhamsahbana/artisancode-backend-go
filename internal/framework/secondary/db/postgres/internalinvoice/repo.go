package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type internalInvoiceRepo struct {
	db *sqlx.DB
}

type Config struct {
	DB *sqlx.DB
}

var _ portsRepo.InternalInvoiceRepository = &internalInvoiceRepo{}

func NewInternalInvoiceRepository(cfg Config) portsRepo.InternalInvoiceRepository {
	return &internalInvoiceRepo{db: cfg.DB}
}
