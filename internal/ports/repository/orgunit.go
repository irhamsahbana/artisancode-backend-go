package repository

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type OrgUnitRepository interface {
	GetOrgUnits(ctx context.Context, filter coreentity.OrgUnitListFilter) ([]coreentity.OrgUnit, int, error)
	GetAllOrgUnitsByCompany(ctx context.Context, tenantID string, companyID string) ([]coreentity.OrgUnit, error)
	GetOrgUnit(ctx context.Context, filter coreentity.OrgUnit) (*coreentity.OrgUnit, error)
	CreateOrgUnit(ctx context.Context, data coreentity.OrgUnit) (*coreentity.OrgUnit, error)
	UpdateOrgUnit(ctx context.Context, data coreentity.OrgUnit) error
	DeleteOrgUnit(ctx context.Context, filter coreentity.OrgUnitDeleteFilter) error
	ExistsOrgUnitByCode(ctx context.Context, tenantID string, code string) (bool, error)
}
