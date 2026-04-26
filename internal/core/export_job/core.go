package core

import (
	corePorts "codebase-app/internal/ports/core"
	integrationPorts "codebase-app/internal/ports/integration"
	portsRepo "codebase-app/internal/ports/secondary/db"
)

type exportJobCore struct {
	repo           portsRepo.ExportJobRepository
	attendanceRepo portsRepo.AttendanceRepository
	storageRepo    portsRepo.StorageRepository
	tx             portsRepo.Transactor
	s3             integrationPorts.StorageContract
	bus            integrationPorts.MessagePublisher
}

var _ corePorts.ExportJobCore = &exportJobCore{}

type Config struct {
	Repo           portsRepo.ExportJobRepository
	AttendanceRepo portsRepo.AttendanceRepository
	StorageRepo    portsRepo.StorageRepository
	Tx             portsRepo.Transactor
	S3             integrationPorts.StorageContract
	Bus            integrationPorts.MessagePublisher
}

func NewExportJobCore(cfg Config) *exportJobCore {
	return &exportJobCore{
		repo:           cfg.Repo,
		attendanceRepo: cfg.AttendanceRepo,
		storageRepo:    cfg.StorageRepo,
		tx:             cfg.Tx,
		s3:             cfg.S3,
		bus:            cfg.Bus,
	}
}
