package repository

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type WorkLocationRepository interface {
	GetWorkLocations(ctx context.Context, filter coreentity.WorkLocationListFilter) ([]coreentity.WorkLocation, int, error)
	GetWorkLocation(ctx context.Context, filter coreentity.WorkLocation) (*coreentity.WorkLocation, error)
	CreateWorkLocation(ctx context.Context, data coreentity.WorkLocation) (*coreentity.WorkLocation, error)
	UpdateWorkLocation(ctx context.Context, data coreentity.WorkLocation) error
	DeleteWorkLocation(ctx context.Context, filter coreentity.WorkLocationDeleteFilter) error
	ExistsWorkLocationByName(ctx context.Context, tenantID string, name string, excludeID *string) (bool, error)
}