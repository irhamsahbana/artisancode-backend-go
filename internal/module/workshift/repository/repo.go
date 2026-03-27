package repository

import (
	portsRepo "codebase-app/internal/ports/repository"

	"github.com/jmoiron/sqlx"
)

type workShiftRepo struct {
	db *sqlx.DB
}

var _ portsRepo.WorkShiftRepository = &workShiftRepo{}

type WorkShiftRepositoryConfig struct {
	DB *sqlx.DB
}

func NewWorkShiftRepository(cfg WorkShiftRepositoryConfig) portsRepo.WorkShiftRepository {
	return &workShiftRepo{db: cfg.DB}
}