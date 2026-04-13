package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"context"
)

type orgUnitCore struct {
	repo portsRepo.OrgUnitRepository
}

type OrgUnitCoreConfig struct {
	Repo portsRepo.OrgUnitRepository
}

var _ corePorts.OrgUnitCore = &orgUnitCore{}

func NewOrgUnitCore(cfg OrgUnitCoreConfig) *orgUnitCore {
	return &orgUnitCore{repo: cfg.Repo}
}

func (c *orgUnitCore) GetOrgUnits(ctx context.Context, filter coreentity.OrgUnitListFilter) ([]coreentity.OrgUnit, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:orgunit:core:GetOrgUnits")
	defer span.End()

	return c.repo.GetOrgUnits(ctx, filter)
}

func (c *orgUnitCore) GetOrgUnit(ctx context.Context, filter coreentity.OrgUnit) (*coreentity.OrgUnit, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:orgunit:core:GetOrgUnit")
	defer span.End()

	return c.repo.GetOrgUnit(ctx, filter)
}

func (c *orgUnitCore) UpdateOrgUnit(ctx context.Context, data coreentity.OrgUnit) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:orgunit:core:UpdateOrgUnit")
	defer span.End()

	// Check if org unit exists before updating
	_, err := c.repo.GetOrgUnit(ctx, coreentity.OrgUnit{
		TenantID: data.TenantID,
		ID:       data.ID,
	})
	if err != nil {
		return err
	}

	return c.repo.UpdateOrgUnit(ctx, data)
}

func (c *orgUnitCore) GetAllOrgUnitsByCompany(ctx context.Context, tenantID string, companyID string) ([]coreentity.OrgUnit, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:orgunit:core:GetAllOrgUnitsByCompany")
	defer span.End()

	return c.repo.GetAllOrgUnitsByCompany(ctx, tenantID, companyID)
}
