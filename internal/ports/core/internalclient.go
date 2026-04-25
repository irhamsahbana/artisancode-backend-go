package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalClientCore interface {
	GetInternalClients(ctx context.Context, filter coreentity.InternalClientListFilter) ([]coreentity.InternalClient, int, error)
}
