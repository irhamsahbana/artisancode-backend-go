package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type internalUserRepo struct {
	db *sqlx.DB
}

type Config struct {
	DB *sqlx.DB
}

var _ portsRepo.InternalUserRepository = &internalUserRepo{}

func NewInternalUserRepository(cfg Config) portsRepo.InternalUserRepository {
	return &internalUserRepo{db: cfg.DB}
}
