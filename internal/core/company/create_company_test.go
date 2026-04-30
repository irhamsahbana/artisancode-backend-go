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

func TestCompanyCore_CreateCompany(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     coreentity.Company
		setup     func(repo *mocks.CompanyRepository)
		want      *coreentity.Company
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.Company{
				TenantID: "tenant-1",
				Name:     "Acme Corp",
				Code:     "ACME",
				Config: coreentity.CompanyConfig{
					Timezone:           "Asia/Makassar",
					DateFormat:         "YYYY-MM-DD",
					TimeFormat:         "HH:mm:ss",
					PreferredLanguage:  "id",
					SupportedLanguages: []string{"id", "en"},
				},
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					ExistsCompanyByCode(mock.Anything, "tenant-1", "ACME", "").
					Return(false, nil)
				repo.EXPECT().
					CreateCompany(mock.Anything, mock.AnythingOfType("coreentity.Company")).
					Return(&coreentity.Company{
						ID:       "company-1",
						TenantID: "tenant-1",
						Name:     "Acme Corp",
						Code:     "ACME",
					}, nil)
			},
			want: &coreentity.Company{
				ID:       "company-1",
				TenantID: "tenant-1",
				Name:     "Acme Corp",
				Code:     "ACME",
			},
		},
		{
			name: "validation error - timezone required",
			input: coreentity.Company{
				TenantID: "tenant-1",
				Name:     "Acme Corp",
				Code:     "ACME",
				Config: coreentity.CompanyConfig{
					Timezone:           "",
					DateFormat:         "YYYY-MM-DD",
					TimeFormat:         "HH:mm:ss",
					PreferredLanguage:  "id",
					SupportedLanguages: []string{"id", "en"},
				},
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					ExistsCompanyByCode(mock.Anything, "tenant-1", "ACME", "").
					Return(false, nil)
			},
			wantError: true,
		},
		{
			name: "validation error - date format required",
			input: coreentity.Company{
				TenantID: "tenant-1",
				Name:     "Acme Corp",
				Code:     "ACME",
				Config: coreentity.CompanyConfig{
					Timezone:           "Asia/Makassar",
					DateFormat:         "",
					TimeFormat:         "HH:mm:ss",
					PreferredLanguage:  "id",
					SupportedLanguages: []string{"id", "en"},
				},
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					ExistsCompanyByCode(mock.Anything, "tenant-1", "ACME", "").
					Return(false, nil)
			},
			wantError: true,
		},
		{
			name: "validation error - time format required",
			input: coreentity.Company{
				TenantID: "tenant-1",
				Name:     "Acme Corp",
				Code:     "ACME",
				Config: coreentity.CompanyConfig{
					Timezone:           "Asia/Makassar",
					DateFormat:         "YYYY-MM-DD",
					TimeFormat:         "",
					PreferredLanguage:  "id",
					SupportedLanguages: []string{"id", "en"},
				},
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					ExistsCompanyByCode(mock.Anything, "tenant-1", "ACME", "").
					Return(false, nil)
			},
			wantError: true,
		},
		{
			name: "validation error - preferred language required",
			input: coreentity.Company{
				TenantID: "tenant-1",
				Name:     "Acme Corp",
				Code:     "ACME",
				Config: coreentity.CompanyConfig{
					Timezone:           "Asia/Makassar",
					DateFormat:         "YYYY-MM-DD",
					TimeFormat:         "HH:mm:ss",
					PreferredLanguage:  "",
					SupportedLanguages: []string{"id", "en"},
				},
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					ExistsCompanyByCode(mock.Anything, "tenant-1", "ACME", "").
					Return(false, nil)
			},
			wantError: true,
		},
		{
			name: "validation error - supported languages required",
			input: coreentity.Company{
				TenantID: "tenant-1",
				Name:     "Acme Corp",
				Code:     "ACME",
				Config: coreentity.CompanyConfig{
					Timezone:           "Asia/Makassar",
					DateFormat:         "YYYY-MM-DD",
					TimeFormat:         "HH:mm:ss",
					PreferredLanguage:  "id",
					SupportedLanguages: []string{},
				},
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					ExistsCompanyByCode(mock.Anything, "tenant-1", "ACME", "").
					Return(false, nil)
			},
			wantError: true,
		},
		{
			name: "validation error - unsupported language",
			input: coreentity.Company{
				TenantID: "tenant-1",
				Name:     "Acme Corp",
				Code:     "ACME",
				Config: coreentity.CompanyConfig{
					Timezone:           "Asia/Makassar",
					DateFormat:         "YYYY-MM-DD",
					TimeFormat:         "HH:mm:ss",
					PreferredLanguage:  "id",
					SupportedLanguages: []string{"id", "fr"},
				},
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					ExistsCompanyByCode(mock.Anything, "tenant-1", "ACME", "").
					Return(false, nil)
			},
			wantError: true,
		},
		{
			name: "validation error - preferred language not in supported",
			input: coreentity.Company{
				TenantID: "tenant-1",
				Name:     "Acme Corp",
				Code:     "ACME",
				Config: coreentity.CompanyConfig{
					Timezone:           "Asia/Makassar",
					DateFormat:         "YYYY-MM-DD",
					TimeFormat:         "HH:mm:ss",
					PreferredLanguage:  "en",
					SupportedLanguages: []string{"id"},
				},
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					ExistsCompanyByCode(mock.Anything, "tenant-1", "ACME", "").
					Return(false, nil)
			},
			wantError: true,
		},
		{
			name: "duplicate code error",
			input: coreentity.Company{
				TenantID: "tenant-1",
				Name:     "Acme Corp",
				Code:     "ACME",
				Config: coreentity.CompanyConfig{
					Timezone:           "Asia/Makassar",
					DateFormat:         "YYYY-MM-DD",
					TimeFormat:         "HH:mm:ss",
					PreferredLanguage:  "id",
					SupportedLanguages: []string{"id", "en"},
				},
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					ExistsCompanyByCode(mock.Anything, "tenant-1", "ACME", "").
					Return(true, nil)
			},
			wantError: true,
		},
		{
			name: "dependency error from repository - exists check",
			input: coreentity.Company{
				TenantID: "tenant-1",
				Name:     "Acme Corp",
				Code:     "ACME",
				Config: coreentity.CompanyConfig{
					Timezone:           "Asia/Makassar",
					DateFormat:         "YYYY-MM-DD",
					TimeFormat:         "HH:mm:ss",
					PreferredLanguage:  "id",
					SupportedLanguages: []string{"id", "en"},
				},
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					ExistsCompanyByCode(mock.Anything, "tenant-1", "ACME", "").
					Return(false, errmsg.NewCustomErrors(500).SetMessage("database error"))
			},
			wantError: true,
		},
		{
			name: "dependency error from repository - create",
			input: coreentity.Company{
				TenantID: "tenant-1",
				Name:     "Acme Corp",
				Code:     "ACME",
				Config: coreentity.CompanyConfig{
					Timezone:           "Asia/Makassar",
					DateFormat:         "YYYY-MM-DD",
					TimeFormat:         "HH:mm:ss",
					PreferredLanguage:  "id",
					SupportedLanguages: []string{"id", "en"},
				},
			},
			setup: func(repo *mocks.CompanyRepository) {
				repo.EXPECT().
					ExistsCompanyByCode(mock.Anything, "tenant-1", "ACME", "").
					Return(false, nil)
				repo.EXPECT().
					CreateCompany(mock.Anything, mock.AnythingOfType("coreentity.Company")).
					Return(nil, errmsg.NewCustomErrors(500).SetMessage("database error"))
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

			got, err := core.CreateCompany(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
