package core

import (
	"codebase-app/internal/integration/tokencache"
	corePorts "codebase-app/internal/ports/core"
	integrationPorts "codebase-app/internal/ports/integration"
	repository "codebase-app/internal/ports/secondary/db"
)

var _ corePorts.UserCore = &userCore{}

type userCore struct {
	repo                 repository.UserRepository
	tx                   repository.Transactor
	tokenCache           tokencache.TokenCacheContract
	bus                  integrationPorts.MessagePublisher
	googleTokenValidator integrationPorts.GoogleIDTokenValidator
}

type Config struct {
	Repo                 repository.UserRepository
	Tx                   repository.Transactor
	TokenCache           tokencache.TokenCacheContract
	Bus                  integrationPorts.MessagePublisher
	GoogleTokenValidator integrationPorts.GoogleIDTokenValidator
}

func NewUserCore(cfg Config) *userCore {
	return &userCore{
		repo:                 cfg.Repo,
		tx:                   cfg.Tx,
		tokenCache:           cfg.TokenCache,
		bus:                  cfg.Bus,
		googleTokenValidator: cfg.GoogleTokenValidator,
	}
}
