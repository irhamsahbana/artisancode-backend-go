package core

import (
	"regexp"

	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
)

var (
	currencyCodePattern = regexp.MustCompile(`^[A-Z]{3}$`)
)

type internalCurrencyCore struct {
	repo portsRepo.InternalCurrencyRepository
}

type Config struct {
	Repo portsRepo.InternalCurrencyRepository
}

var _ corePorts.InternalCurrencyCore = &internalCurrencyCore{}

func NewInternalCurrencyCore(cfg Config) *internalCurrencyCore {
	return &internalCurrencyCore{repo: cfg.Repo}
}
