package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/internal/ports/secondary/integration"
)

type storageCore struct {
	s3   integration.StorageContract
	repo portsRepo.StorageRepository
}

var _ corePorts.StorageCore = &storageCore{}

func NewStorageCore(s3 integration.StorageContract, repo portsRepo.StorageRepository) *storageCore {
	return &storageCore{
		s3:   s3,
		repo: repo,
	}
}
