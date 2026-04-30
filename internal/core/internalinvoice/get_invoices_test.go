package core

import (
	"errors"
	"testing"

	"codebase-app/internal/entity/coreentity"
	dbmocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestInternalInvoiceCore_GetInvoices(t *testing.T) {
	ctx := internalInvoiceContext()

	tests := []struct {
		name      string
		setup     func(repo *dbmocks.InternalInvoiceRepository)
		want      []coreentity.InternalInvoice
		wantCount int
		wantError bool
	}{
		{
			name: "success uses tenant from context",
			setup: func(repo *dbmocks.InternalInvoiceRepository) {
				repo.EXPECT().
					GetInvoices(mock.Anything, coreentity.InternalCommerceListFilter{TenantID: "tenant-1"}).
					Return([]coreentity.InternalInvoice{{ID: "invoice-1"}}, 1, nil)
			},
			want:      []coreentity.InternalInvoice{{ID: "invoice-1"}},
			wantCount: 1,
		},
		{
			name: "repository error",
			setup: func(repo *dbmocks.InternalInvoiceRepository) {
				repo.EXPECT().
					GetInvoices(mock.Anything, coreentity.InternalCommerceListFilter{TenantID: "tenant-1"}).
					Return(nil, 0, errors.New("repository failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbmocks.NewInternalInvoiceRepository(t)
			tt.setup(repo)

			core := NewInternalInvoiceCore(Config{InvoiceRepo: repo})
			got, count, err := core.GetInvoices(ctx, coreentity.InternalCommerceListFilter{})

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
