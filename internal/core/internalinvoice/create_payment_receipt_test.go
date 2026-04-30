package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreatePaymentReceipt(t *testing.T) {
	ctx := context.Background()
	input := coreentity.InternalPaymentReceipt{
		UserCtx:           common.UserContext{TenantID: "tenant-1", UserID: "user-1"},
		InternalInvoiceID: "invoice-1",
		AmountReceived:    decimal.NewFromInt(100),
	}

	tests := []struct {
		name      string
		input     coreentity.InternalPaymentReceipt
		setup     func(invoiceRepo *dbMocks.InternalInvoiceRepository, orderRepo *dbMocks.InternalOrderRepository, tx *dbMocks.Transactor)
		wantError bool
	}{
		{
			name:  "creates pending manual receipt without marking paid",
			input: input,
			setup: func(invoiceRepo *dbMocks.InternalInvoiceRepository, orderRepo *dbMocks.InternalOrderRepository, tx *dbMocks.Transactor) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				invoiceRepo.EXPECT().
					GetInvoice(mock.Anything, "tenant-1", "invoice-1").
					Return(&coreentity.InternalInvoice{
						ID:                "invoice-1",
						AmountOutstanding: decimal.NewFromInt(100),
					}, nil)
				invoiceRepo.EXPECT().
					CreatePaymentReceipt(mock.Anything, mock.MatchedBy(func(data coreentity.InternalPaymentReceipt) bool {
						return data.Status == coreentity.PaymentReceiptStatusPendingVerification &&
							data.SourceType == coreentity.PaymentProviderManual &&
							data.VerifiedAt == nil &&
							data.VerifiedByUserID == nil
					})).
					Return(&coreentity.InternalPaymentReceipt{ID: "receipt-1"}, nil)
				orderRepo.EXPECT().
					GetOrderByInvoiceID(mock.Anything, "tenant-1", "invoice-1").
					Return(&coreentity.InternalOrder{ID: "order-1"}, nil)
			},
		},
		{
			name: "accepts full amount and marks invoice/order paid",
			input: coreentity.InternalPaymentReceipt{
				UserCtx:           common.UserContext{TenantID: "tenant-1", UserID: "user-1"},
				InternalInvoiceID: "invoice-1",
				AmountReceived:    decimal.NewFromInt(100),
				Status:            coreentity.PaymentReceiptStatusAccepted,
			},
			setup: func(invoiceRepo *dbMocks.InternalInvoiceRepository, orderRepo *dbMocks.InternalOrderRepository, tx *dbMocks.Transactor) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				invoiceRepo.EXPECT().
					GetInvoice(mock.Anything, "tenant-1", "invoice-1").
					Return(&coreentity.InternalInvoice{
						ID:                "invoice-1",
						AmountOutstanding: decimal.NewFromInt(100),
					}, nil)
				invoiceRepo.EXPECT().
					CreatePaymentReceipt(mock.Anything, mock.MatchedBy(func(data coreentity.InternalPaymentReceipt) bool {
						return data.Status == coreentity.PaymentReceiptStatusAccepted &&
							data.VerifiedAt != nil &&
							data.VerifiedByUserID != nil &&
							*data.VerifiedByUserID == "user-1"
					})).
					Return(&coreentity.InternalPaymentReceipt{ID: "receipt-1"}, nil)
				orderRepo.EXPECT().
					GetOrderByInvoiceID(mock.Anything, "tenant-1", "invoice-1").
					Return(&coreentity.InternalOrder{ID: "order-1"}, nil)
				invoiceRepo.EXPECT().
					MarkInvoicePaid(mock.Anything, "tenant-1", "invoice-1", "100").
					Return(&coreentity.InternalInvoice{ID: "invoice-1", Status: coreentity.InvoiceStatusPaid}, nil)
				orderRepo.EXPECT().
					MarkOrderPaid(mock.Anything, "tenant-1", "order-1").
					Return(&coreentity.InternalOrder{ID: "order-1", Status: coreentity.OrderStatusPaid}, nil)
			},
		},
		{
			name: "rejects accepted receipt below outstanding amount",
			input: coreentity.InternalPaymentReceipt{
				UserCtx:           common.UserContext{TenantID: "tenant-1", UserID: "user-1"},
				InternalInvoiceID: "invoice-1",
				AmountReceived:    decimal.NewFromInt(99),
				Status:            coreentity.PaymentReceiptStatusAccepted,
			},
			setup: func(invoiceRepo *dbMocks.InternalInvoiceRepository, orderRepo *dbMocks.InternalOrderRepository, tx *dbMocks.Transactor) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				invoiceRepo.EXPECT().
					GetInvoice(mock.Anything, "tenant-1", "invoice-1").
					Return(&coreentity.InternalInvoice{
						ID:                "invoice-1",
						AmountOutstanding: decimal.NewFromInt(100),
					}, nil)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoiceRepo := dbMocks.NewInternalInvoiceRepository(t)
			orderRepo := dbMocks.NewInternalOrderRepository(t)
			tx := dbMocks.NewTransactor(t)
			if tt.setup != nil {
				tt.setup(invoiceRepo, orderRepo, tx)
			}
			core := NewInternalInvoiceCore(Config{InvoiceRepo: invoiceRepo, OrderRepo: orderRepo, Tx: tx})

			got, err := core.CreatePaymentReceipt(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, "receipt-1", got.Receipt.ID)
		})
	}
}
