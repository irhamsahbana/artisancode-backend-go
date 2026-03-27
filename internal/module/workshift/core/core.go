package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/repository"
)

type workShiftCore struct {
	repo portsRepo.WorkShiftRepository
}

type WorkShiftCoreConfig struct {
	Repo portsRepo.WorkShiftRepository
}

var _ corePorts.WorkShiftCore = &workShiftCore{}

func NewWorkShiftCore(cfg WorkShiftCoreConfig) *workShiftCore {
	return &workShiftCore{repo: cfg.Repo}
}