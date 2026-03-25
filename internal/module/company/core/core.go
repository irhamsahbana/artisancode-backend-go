package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/repository"
	"context"
)

type companyCore struct {
	repo portsRepo.CompanyRepository
}

type CompanyCoreConfig struct {
	Repo portsRepo.CompanyRepository
}

var _ corePorts.CompanyCore = &companyCore{}

func NewCompanyCore(cfg CompanyCoreConfig) *companyCore {
	return &companyCore{repo: cfg.Repo}
}

func (c *companyCore) GetCompanies(ctx context.Context, filter coreentity.CompanyListFilter) ([]coreentity.Company, int, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetCompanies")
	defer span.End()

	return c.repo.GetCompanies(ctx, filter)
}

func (c *companyCore) GetCompany(ctx context.Context, filter coreentity.Company) (*coreentity.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetCompany")
	defer span.End()

	return c.repo.GetCompany(ctx, filter)
}

func (c *companyCore) CreateCompany(ctx context.Context, data coreentity.Company) (*coreentity.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "core.CreateCompany")
	defer span.End()

	return c.repo.CreateCompany(ctx, data)
}

func (c *companyCore) UpdateCompany(ctx context.Context, data coreentity.Company) error {
	ctx, span := tracing.StartSpan(ctx, "core.UpdateCompany")
	defer span.End()

	return c.repo.UpdateCompany(ctx, data)
}

func (c *companyCore) DeleteCompany(ctx context.Context, filter coreentity.CompanyDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "core.DeleteCompany")
	defer span.End()

	return c.repo.DeleteCompany(ctx, filter)
}
