package repository

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type CompanyRepository interface {
	GetCompanies(ctx context.Context, filter coreentity.CompanyListFilter) ([]coreentity.Company, int, error)
	GetCompany(ctx context.Context, filter coreentity.Company) (*coreentity.Company, error)
	CreateCompany(ctx context.Context, data coreentity.Company) (*coreentity.Company, error)
	UpdateCompany(ctx context.Context, data coreentity.Company) error
	DeleteCompany(ctx context.Context, filter coreentity.CompanyDeleteFilter) error
}
