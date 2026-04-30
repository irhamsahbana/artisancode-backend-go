package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetQuotations(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{TenantID: "tenant-1"})
	repo := dbMocks.NewInternalQuotationRepository(t)
	repo.EXPECT().
		GetQuotations(mock.Anything, coreentity.InternalCommerceListFilter{TenantID: "tenant-1"}).
		Return([]coreentity.InternalQuotation{{ID: "quote-1"}}, 1, nil)
	core := NewInternalQuotationCore(Config{QuotationRepo: repo})

	got, total, err := core.GetQuotations(ctx, coreentity.InternalCommerceListFilter{})

	require.NoError(t, err)
	require.Equal(t, []coreentity.InternalQuotation{{ID: "quote-1"}}, got)
	require.Equal(t, 1, total)
}
