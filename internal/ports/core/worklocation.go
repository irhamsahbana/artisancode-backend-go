package core

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type WorkLocationCore interface {
	GetWorkLocations(ctx context.Context, filter coreentity.WorkLocationListFilter) ([]coreentity.WorkLocation, int, error)
	GetWorkLocation(ctx context.Context, filter coreentity.WorkLocation) (*coreentity.WorkLocation, error)
	CreateWorkLocation(ctx context.Context, data coreentity.WorkLocation) (*coreentity.WorkLocation, error)
	UpdateWorkLocation(ctx context.Context, data coreentity.WorkLocation) error
	DeleteWorkLocation(ctx context.Context, filter coreentity.WorkLocationDeleteFilter) error
}
