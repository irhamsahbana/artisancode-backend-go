package core

import (
	corePorts "codebase-app/internal/ports/core"
	integrationPorts "codebase-app/internal/ports/integration"
	portsRepo "codebase-app/internal/ports/repository"
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
