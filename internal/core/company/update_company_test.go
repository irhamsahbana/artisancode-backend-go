package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/ports/secondary/db/mocks"
	"codebase-app/pkg/errmsg"
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCompanyCore_UpdateCompany(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     coreentity.Company
		setup     func(repo *mocks.CompanyRepository)
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.Company{
				ID:       "company-1",
				TenantID: "tenant-1",
				Name:     "Updated Corp",
				Code:     "UPDATED",
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					GetCompany(mock.Anything, coreentity.Company{
						TenantID: "tenant-1",
						ID:       "company-1",
					}).
					Return(&coreentity.Company{
						ID:       "company-1",
						TenantID: "tenant-1",
						Name:     "Acme Corp",
						Code:     "ACME",
					}, nil)
				repo.EXPECT().
					UpdateCompany(mock.Anything, mock.AnythingOfType("coreentity.Company")).
					Return(nil)
			},
		},
		{
			name: "company not found",
			input: coreentity.Company{
				ID:       "company-1",
				TenantID: "tenant-1",
				Name:     "Updated Corp",
				Code:     "UPDATED",
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					GetCompany(mock.Anything, coreentity.Company{
						TenantID: "tenant-1",
						ID:       "company-1",
					}).
					Return(nil, nil)
			},
			wantError: true,
		},
		{
			name: "dependency error from repository - get company",
			input: coreentity.Company{
				ID:       "company-1",
				TenantID: "tenant-1",
				Name:     "Updated Corp",
				Code:     "UPDATED",
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					GetCompany(mock.Anything, coreentity.Company{
						TenantID: "tenant-1",
						ID:       "company-1",
					}).
					Return(nil, errmsg.NewCustomErrors(500).SetMessage("database error"))
			},
			wantError: true,
		},
		{
			name: "dependency error from repository - update",
			input: coreentity.Company{
				ID:       "company-1",
				TenantID: "tenant-1",
				Name:     "Updated Corp",
				Code:     "UPDATED",
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					GetCompany(mock.Anything, coreentity.Company{
						TenantID: "tenant-1",
						ID:       "company-1",
					}).
					Return(&coreentity.Company{
						ID:       "company-1",
						TenantID: "tenant-1",
						Name:     "Acme Corp",
						Code:     "ACME",
					}, nil)
				repo.EXPECT().
					UpdateCompany(mock.Anything, mock.AnythingOfType("coreentity.Company")).
					Return(errmsg.NewCustomErrors(500).SetMessage("database error"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewCompanyRepository(t)

			if tt.setup != nil {
				tt.setup(repo)
			}

			core := NewCompanyCore(Config{
				Repo: repo,
			})

			err := core.UpdateCompany(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
