package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
)

type exportJobCore struct {
	repo           portsRepo.ExportJobRepository
	attendanceRepo portsRepo.AttendanceRepository
	storageRepo    portsRepo.StorageRepository
	s3             integrationPorts.StorageContract
	bus            integrationPorts.MessagePublisher
}

var _ corePorts.ExportJobCore = &exportJobCore{}

type ExportJobCoreConfig struct {
	Repo           portsRepo.ExportJobRepository
	AttendanceRepo portsRepo.AttendanceRepository
	StorageRepo    portsRepo.StorageRepository
	S3             integrationPorts.StorageContract
	Bus            integrationPorts.MessagePublisher
}

func NewExportJobCore(cfg ExportJobCoreConfig) *exportJobCore {
	return &exportJobCore{
		repo:           cfg.Repo,
		attendanceRepo: cfg.AttendanceRepo,
		storageRepo:    cfg.StorageRepo,
		s3:             cfg.S3,
		bus:            cfg.Bus,
	}
}
