package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalClientRepository interface {
	GetInternalClients(ctx context.Context, filter coreentity.InternalClientListFilter) ([]coreentity.InternalClient, int, error)
}
