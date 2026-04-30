package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/ports/secondary/db/mocks"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCompanyCore_DeleteCompany(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     coreentity.CompanyDeleteFilter
		setup     func(repo *mocks.CompanyRepository)
		wantError bool
	}{
		{
			name: "always returns 403 error - company cannot be deleted",
			input: coreentity.CompanyDeleteFilter{
				ID:       "company-1",
				TenantID: "tenant-1",
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

			err := core.DeleteCompany(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
