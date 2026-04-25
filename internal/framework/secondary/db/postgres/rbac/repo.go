package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type rbacRepo struct {
	db *sqlx.DB
}

var _ portsRepo.RbacRepository = &rbacRepo{}

type RbacRepositoryConfig struct {
	DB *sqlx.DB
}

func NewRbacRepository(cfg RbacRepositoryConfig) portsRepo.RbacRepository {
	return &rbacRepo{
		db: cfg.DB,
	}
}
