package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type workShiftRepo struct {
	db *sqlx.DB
}

var _ portsRepo.WorkShiftRepository = &workShiftRepo{}

type Config struct {
	DB *sqlx.DB
}

func NewWorkShiftRepository(cfg Config) portsRepo.WorkShiftRepository {
	return &workShiftRepo{db: cfg.DB}
}
