package core

import (
	corePorts "codebase-app/internal/ports/core"
	"codebase-app/internal/ports/integration"
	portsRepo "codebase-app/internal/ports/secondary/db"
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
