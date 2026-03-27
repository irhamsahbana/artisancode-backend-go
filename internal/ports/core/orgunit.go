package core

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type OrgUnitCore interface {
	GetOrgUnits(ctx context.Context, filter coreentity.OrgUnitListFilter) ([]coreentity.OrgUnit, int, error)
	GetOrgUnitTree(ctx context.Context, tenantID string, companyID string) ([]coreentity.OrgUnitTreeNode, error)
	GetOrgUnit(ctx context.Context, filter coreentity.OrgUnit) (*coreentity.OrgUnit, error)
	CreateOrgUnit(ctx context.Context, data coreentity.OrgUnit) (*coreentity.OrgUnit, error)
	UpdateOrgUnit(ctx context.Context, data coreentity.OrgUnit) error
	DeleteOrgUnit(ctx context.Context, filter coreentity.OrgUnitDeleteFilter) error
}
