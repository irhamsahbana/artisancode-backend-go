package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type companyRepo struct {
	db *sqlx.DB
}

var _ portsRepo.CompanyRepository = &companyRepo{}

type Config struct {
	DB *sqlx.DB
}

func NewCompanyRepository(cfg Config) portsRepo.CompanyRepository {
	return &companyRepo{db: cfg.DB}
}
