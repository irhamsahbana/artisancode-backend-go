package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type internalProductRepo struct {
	db *sqlx.DB
}

type Config struct {
	DB *sqlx.DB
}

var _ portsRepo.InternalProductRepository = &internalProductRepo{}

func NewInternalProductRepository(cfg Config) portsRepo.InternalProductRepository {
	return &internalProductRepo{db: cfg.DB}
}
