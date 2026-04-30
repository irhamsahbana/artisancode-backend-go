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

func TestGetOrder(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{TenantID: "tenant-1"})
	orderRepo := dbMocks.NewInternalOrderRepository(t)
	invoiceRepo := dbMocks.NewInternalInvoiceRepository(t)
	orderRepo.EXPECT().
		GetOrder(mock.Anything, "tenant-1", "order-1").
		Return(&coreentity.InternalOrder{ID: "order-1"}, nil)
	invoiceRepo.EXPECT().
		GetInvoiceByOrderID(mock.Anything, "tenant-1", "order-1").
		Return(nil, assertErr{})
	core := NewInternalOrderCore(Config{OrderRepo: orderRepo, InvoiceRepo: invoiceRepo})

	got, err := core.GetOrder(ctx, "order-1")

	require.NoError(t, err)
	require.Equal(t, "order-1", got.Order.ID)
	require.Nil(t, got.Invoice)
}

type assertErr struct{}

func (assertErr) Error() string {
	return "ignored"
}
