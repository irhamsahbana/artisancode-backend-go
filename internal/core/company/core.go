package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/pkg/errmsg"
	"context"
)

type companyCore struct {
	repo portsRepo.CompanyRepository
}

type Config struct {
	Repo portsRepo.CompanyRepository
}

var _ corePorts.CompanyCore = &companyCore{}

func NewCompanyCore(cfg Config) *companyCore {
	return &companyCore{repo: cfg.Repo}
}

func (c *companyCore) GetCompanies(ctx context.Context, filter coreentity.CompanyListFilter) ([]coreentity.Company, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:company:core:GetCompanies")
	defer span.End()

	return c.repo.GetCompanies(ctx, filter)
}

func (c *companyCore) GetCompany(ctx context.Context, filter coreentity.Company) (*coreentity.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:company:core:GetCompany")
	defer span.End()

	return c.repo.GetCompany(ctx, filter)
}

func (c *companyCore) DeleteCompany(ctx context.Context, filter coreentity.CompanyDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:company:core:DeleteCompany")
	defer span.End()

	return errmsg.NewCustomErrors(403).SetMessage("Company cannot be deleted")
}
