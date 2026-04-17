package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type messageQueueRepo struct {
	db *sqlx.DB
}

var _ portsRepo.MessageQueueRepository = &messageQueueRepo{}

type MessageQueueRepositoryConfig struct {
	DB *sqlx.DB
}

func NewMessageQueueRepository(cfg MessageQueueRepositoryConfig) portsRepo.MessageQueueRepository {
	return &messageQueueRepo{db: cfg.DB}
}
