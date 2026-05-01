package core

import (
	"regexp"

	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
)

var (
	internalCodePattern     = regexp.MustCompile(`^[A-Z0-9_-]+$`)
	internalCurrencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)
)

type internalProductCore struct {
	repo portsRepo.InternalProductRepository
}

type Config struct {
	Repo portsRepo.InternalProductRepository
}

var _ corePorts.InternalProductCore = &internalProductCore{}

func NewInternalProductCore(cfg Config) *internalProductCore {
	return &internalProductCore{repo: cfg.Repo}
}
