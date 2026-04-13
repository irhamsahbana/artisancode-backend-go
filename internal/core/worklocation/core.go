package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
)

type workLocationCore struct {
	repo        portsRepo.WorkLocationRepository
	orgUnitRepo portsRepo.OrgUnitRepository
}

type WorkLocationCoreConfig struct {
	Repo        portsRepo.WorkLocationRepository
	OrgUnitRepo portsRepo.OrgUnitRepository
}

var _ corePorts.WorkLocationCore = &workLocationCore{}

func NewWorkLocationCore(cfg WorkLocationCoreConfig) *workLocationCore {
	return &workLocationCore{
		repo:        cfg.Repo,
		orgUnitRepo: cfg.OrgUnitRepo,
	}
}
