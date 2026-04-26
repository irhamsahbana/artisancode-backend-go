package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type orgUnitRepo struct {
	db *sqlx.DB
}

var _ portsRepo.OrgUnitRepository = &orgUnitRepo{}

type Config struct {
	DB *sqlx.DB
}

func NewOrgUnitRepository(cfg Config) portsRepo.OrgUnitRepository {
	return &orgUnitRepo{db: cfg.DB}
}
