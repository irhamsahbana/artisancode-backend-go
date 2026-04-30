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

func TestCompanyCore_GetCompanies(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		filter    coreentity.CompanyListFilter
		setup     func(repo *mocks.CompanyRepository)
		want      []coreentity.Company
		wantCount int
		wantError bool
	}{
		{
			name: "success",
			filter: coreentity.CompanyListFilter{
				TenantID: "tenant-1",
				Q:        "acme",
				Page:     1,
				Paginate: 10,
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					GetCompanies(mock.Anything, coreentity.CompanyListFilter{
						TenantID: "tenant-1",
						Q:        "acme",
						Page:     1,
						Paginate: 10,
					}).
					Return([]coreentity.Company{
						{
							ID:       "company-1",
							TenantID: "tenant-1",
							Name:     "Acme Corp",
							Code:     "ACME",
						},
						{
							ID:       "company-2",
							TenantID: "tenant-1",
							Name:     "Acme Subsidiary",
							Code:     "ACME-SUB",
						},
					}, 2, nil)
			},
			want: []coreentity.Company{
				{
					ID:       "company-1",
					TenantID: "tenant-1",
					Name:     "Acme Corp",
					Code:     "ACME",
				},
				{
					ID:       "company-2",
					TenantID: "tenant-1",
					Name:     "Acme Subsidiary",
					Code:     "ACME-SUB",
				},
			},
			wantCount: 2,
		},
		{
			name: "dependency error from repository",
			filter: coreentity.CompanyListFilter{
				TenantID: "tenant-1",
				Q:        "acme",
				Page:     1,
				Paginate: 10,
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					GetCompanies(mock.Anything, coreentity.CompanyListFilter{
						TenantID: "tenant-1",
						Q:        "acme",
						Page:     1,
						Paginate: 10,
					}).
					Return(nil, 0, errmsg.NewCustomErrors(500).SetMessage("database error"))
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

			got, count, err := core.GetCompanies(ctx, tt.filter)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantCount, count)
		})
	}
}
