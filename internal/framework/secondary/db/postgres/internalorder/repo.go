package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type internalOrderRepo struct {
	db *sqlx.DB
}

type Config struct {
	DB *sqlx.DB
}

var _ portsRepo.InternalOrderRepository = &internalOrderRepo{}

func NewInternalOrderRepository(cfg Config) portsRepo.InternalOrderRepository {
	return &internalOrderRepo{db: cfg.DB}
}
