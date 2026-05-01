package core

import (
	"codebase-app/internal/integration/tokencache"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
)

type internalUserCore struct {
	repo       portsRepo.InternalUserRepository
	tokenCache tokencache.TokenCacheContract
}

type Config struct {
	Repo       portsRepo.InternalUserRepository
	TokenCache tokencache.TokenCacheContract
}

var _ corePorts.InternalUserCore = &internalUserCore{}

func NewInternalUserCore(cfg Config) *internalUserCore {
	return &internalUserCore{
		repo:       cfg.Repo,
		tokenCache: cfg.TokenCache,
	}
}
