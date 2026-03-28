package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/repository"
)

var _ corePorts.MeCore = &meCore{}

type meCore struct {
	repo portsRepo.MeRepository
}

type MeCoreConfig struct {
	Repo portsRepo.MeRepository
}

func NewMeCore(cfg MeCoreConfig) *meCore {
	return &meCore{
		repo: cfg.Repo,
	}
}
