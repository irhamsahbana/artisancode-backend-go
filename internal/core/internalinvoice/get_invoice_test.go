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

func TestGetInvoice(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{TenantID: "tenant-1"})
	quoteID := "quote-1"
	invoiceRepo := dbMocks.NewInternalInvoiceRepository(t)
	orderRepo := dbMocks.NewInternalOrderRepository(t)
	quotationRepo := dbMocks.NewInternalQuotationRepository(t)
	invoiceRepo.EXPECT().
		GetInvoice(mock.Anything, "tenant-1", "invoice-1").
		Return(&coreentity.InternalInvoice{ID: "invoice-1"}, nil)
	orderRepo.EXPECT().
		GetOrderByInvoiceID(mock.Anything, "tenant-1", "invoice-1").
		Return(&coreentity.InternalOrder{
			ID:                "order-1",
			SourceType:        coreentity.OrderSourceTypeQuotation,
			SourceReferenceID: &quoteID,
		}, nil)
	quotationRepo.EXPECT().
		GetQuotation(mock.Anything, "tenant-1", "quote-1").
		Return(&coreentity.InternalQuotation{ID: "quote-1"}, nil)
	invoiceRepo.EXPECT().
		GetLatestPaymentAttempt(mock.Anything, "tenant-1", "invoice-1").
		Return(&coreentity.InternalPaymentAttempt{ID: "attempt-1"}, nil)
	core := NewInternalInvoiceCore(Config{
		InvoiceRepo:   invoiceRepo,
		OrderRepo:     orderRepo,
		QuotationRepo: quotationRepo,
	})

	got, err := core.GetInvoice(ctx, "invoice-1")

	require.NoError(t, err)
	require.Equal(t, "invoice-1", got.Invoice.ID)
	require.Equal(t, "order-1", got.Order.ID)
	require.Equal(t, "quote-1", got.Quotation.ID)
	require.Equal(t, "attempt-1", got.LatestPaymentAttempt.ID)
}
