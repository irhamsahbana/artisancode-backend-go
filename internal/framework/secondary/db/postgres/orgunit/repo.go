package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type orgUnitRepo struct {
	db *sqlx.DB
}

var _ portsRepo.OrgUnitRepository = &orgUnitRepo{}

type OrgUnitRepositoryConfig struct {
	DB *sqlx.DB
}

func NewOrgUnitRepository(cfg OrgUnitRepositoryConfig) portsRepo.OrgUnitRepository {
	return &orgUnitRepo{db: cfg.DB}
}
