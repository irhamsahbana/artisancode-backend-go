package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type employeeRepo struct {
	db *sqlx.DB
}

var _ portsRepo.EmployeeRepository = &employeeRepo{}

type EmployeeRepositoryConfig struct {
	DB *sqlx.DB
}

func NewEmployeeRepository(cfg EmployeeRepositoryConfig) portsRepo.EmployeeRepository {
	return &employeeRepo{db: cfg.DB}
}
