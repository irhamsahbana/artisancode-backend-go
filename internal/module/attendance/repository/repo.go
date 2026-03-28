package repository

import (
	portsRepo "codebase-app/internal/ports/repository"

	"github.com/jmoiron/sqlx"
)

type attendanceRepo struct {
	db *sqlx.DB
}

var _ portsRepo.AttendanceRepository = &attendanceRepo{}

type AttendanceRepositoryConfig struct {
	DB *sqlx.DB
}

func NewAttendanceRepository(cfg AttendanceRepositoryConfig) portsRepo.AttendanceRepository {
	return &attendanceRepo{db: cfg.DB}
}
