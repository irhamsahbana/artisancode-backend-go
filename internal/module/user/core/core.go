package core

import (
	"codebase-app/internal/integration/tokencache"
	corePorts "codebase-app/internal/ports/core"
	"codebase-app/internal/ports/repository"
)

var _ corePorts.UserCore = &userCore{}

type userCore struct {
	repo       repository.UserRepository
	tokenCache tokencache.TokenCacheContract
}

type UserCoreConfig struct {
	Repo       repository.UserRepository
	TokenCache tokencache.TokenCacheContract
}

func NewUserCore(cfg UserCoreConfig) *userCore {
	return &userCore{
		repo:       cfg.Repo,
		tokenCache: cfg.TokenCache,
	}
}
