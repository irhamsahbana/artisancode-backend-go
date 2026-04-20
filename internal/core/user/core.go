package core

import (
	"codebase-app/internal/integration/tokencache"
	corePorts "codebase-app/internal/ports/core"
	"codebase-app/internal/ports/secondary/db"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
)

var _ corePorts.UserCore = &userCore{}

type userCore struct {
	repo       repository.UserRepository
	tokenCache tokencache.TokenCacheContract
	bus        integrationPorts.MessagePublisher
}

type UserCoreConfig struct {
	Repo       repository.UserRepository
	TokenCache tokencache.TokenCacheContract
	Bus        integrationPorts.MessagePublisher
}

func NewUserCore(cfg UserCoreConfig) *userCore {
	return &userCore{
		repo:       cfg.Repo,
		tokenCache: cfg.TokenCache,
		bus:        cfg.Bus,
	}
}
