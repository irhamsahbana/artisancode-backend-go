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

func TestGetPaymentAttempts(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{TenantID: "tenant-1"})
	repo := dbMocks.NewInternalInvoiceRepository(t)
	repo.EXPECT().
		GetPaymentAttempts(mock.Anything, "tenant-1", "invoice-1").
		Return([]coreentity.InternalPaymentAttempt{{ID: "attempt-1"}}, nil)
	core := NewInternalInvoiceCore(Config{InvoiceRepo: repo})

	got, err := core.GetPaymentAttempts(ctx, "invoice-1")

	require.NoError(t, err)
	require.Equal(t, []coreentity.InternalPaymentAttempt{{ID: "attempt-1"}}, got)
}
