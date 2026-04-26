package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
)

var _ corePorts.MeCore = &meCore{}

type meCore struct {
	repo portsRepo.MeRepository
}

type Config struct {
	Repo portsRepo.MeRepository
}

func NewMeCore(cfg Config) *meCore {
	return &meCore{
		repo: cfg.Repo,
	}
}
