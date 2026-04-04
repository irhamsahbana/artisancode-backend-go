package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/repository"
)

var _ corePorts.AttendanceCore = &attendanceCore{}

type attendanceCore struct {
	repo        portsRepo.AttendanceRepository
	companyRepo portsRepo.CompanyRepository
	storageRepo portsRepo.StorageRepository
}

type AttendanceCoreConfig struct {
	Repo        portsRepo.AttendanceRepository
	CompanyRepo portsRepo.CompanyRepository
	StorageRepo portsRepo.StorageRepository
}

func NewAttendanceCore(cfg AttendanceCoreConfig) *attendanceCore {
	return &attendanceCore{
		repo:        cfg.Repo,
		companyRepo: cfg.CompanyRepo,
		storageRepo: cfg.StorageRepo,
	}
}
