package core

import portsRepo "codebase-app/internal/ports/secondary/db"

type messageQueueCore struct {
	repo portsRepo.MessageQueueRepository
}

type MessageQueueCoreConfig struct {
	Repo portsRepo.MessageQueueRepository
}

func NewMessageQueueCore(cfg MessageQueueCoreConfig) *messageQueueCore {
	return &messageQueueCore{
		repo: cfg.Repo,
	}
}
