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

func TestGetQuotation(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{TenantID: "tenant-1"})
	repo := dbMocks.NewInternalQuotationRepository(t)
	repo.EXPECT().
		GetQuotation(mock.Anything, "tenant-1", "quote-1").
		Return(&coreentity.InternalQuotation{ID: "quote-1"}, nil)
	core := NewInternalQuotationCore(Config{QuotationRepo: repo})

	got, err := core.GetQuotation(ctx, "quote-1")

	require.NoError(t, err)
	require.Equal(t, "quote-1", got.ID)
}
