package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
)

type workShiftCore struct {
	repo portsRepo.WorkShiftRepository
}

type Config struct {
	Repo portsRepo.WorkShiftRepository
}

var _ corePorts.WorkShiftCore = &workShiftCore{}

func NewWorkShiftCore(cfg Config) *workShiftCore {
	return &workShiftCore{repo: cfg.Repo}
}
