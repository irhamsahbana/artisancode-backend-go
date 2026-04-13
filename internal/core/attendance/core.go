package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
)

var _ corePorts.AttendanceCore = &attendanceCore{}

type attendanceCore struct {
	repo        portsRepo.AttendanceRepository
	companyRepo portsRepo.CompanyRepository
	storageRepo portsRepo.StorageRepository
	s3          integrationPorts.StorageContract
}

type AttendanceCoreConfig struct {
	Repo        portsRepo.AttendanceRepository
	CompanyRepo portsRepo.CompanyRepository
	StorageRepo portsRepo.StorageRepository
	S3          integrationPorts.StorageContract
}

func NewAttendanceCore(cfg AttendanceCoreConfig) *attendanceCore {
	return &attendanceCore{
		repo:        cfg.Repo,
		companyRepo: cfg.CompanyRepo,
		storageRepo: cfg.StorageRepo,
		s3:          cfg.S3,
	}
}
