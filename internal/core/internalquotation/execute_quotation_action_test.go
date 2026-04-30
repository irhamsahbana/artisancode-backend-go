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

func TestExecuteQuotationAction(t *testing.T) {
	ctx := context.Background()
	input := coreentity.InternalQuotationActionInput{
		UserCtx: common.UserContext{TenantID: "tenant-1"},
		ID:      "quote-1",
		Action:  " " + coreentity.ActionApproveQuotation + " ",
	}

	tests := []struct {
		name      string
		input     coreentity.InternalQuotationActionInput
		setup     func(quoteRepo *dbMocks.InternalQuotationRepository, orderRepo *dbMocks.InternalOrderRepository, invoiceRepo *dbMocks.InternalInvoiceRepository, tx *dbMocks.Transactor)
		wantError bool
	}{
		{
			name:  "approves draft quotation",
			input: input,
			setup: func(quoteRepo *dbMocks.InternalQuotationRepository, orderRepo *dbMocks.InternalOrderRepository, invoiceRepo *dbMocks.InternalInvoiceRepository, tx *dbMocks.Transactor) {
				quoteRepo.EXPECT().
					GetQuotation(mock.Anything, "tenant-1", "quote-1").
					Return(&coreentity.InternalQuotation{ID: "quote-1", Status: coreentity.QuotationStatusDraft}, nil)
				quoteRepo.EXPECT().
					ApproveQuotation(mock.Anything, "tenant-1", "quote-1").
					Return(&coreentity.InternalQuotation{ID: "quote-1", Status: coreentity.QuotationStatusApproved}, nil)
			},
		},
		{
			name: "rejects unsupported action",
			input: coreentity.InternalQuotationActionInput{
				UserCtx: common.UserContext{TenantID: "tenant-1"},
				ID:      "quote-1",
				Action:  "unsupported",
			},
			wantError: true,
		},
		{
			name: "rejects invalid approve status",
			input: coreentity.InternalQuotationActionInput{
				UserCtx: common.UserContext{TenantID: "tenant-1"},
				ID:      "quote-1",
				Action:  coreentity.ActionApproveQuotation,
			},
			setup: func(quoteRepo *dbMocks.InternalQuotationRepository, orderRepo *dbMocks.InternalOrderRepository, invoiceRepo *dbMocks.InternalInvoiceRepository, tx *dbMocks.Transactor) {
				quoteRepo.EXPECT().
					GetQuotation(mock.Anything, "tenant-1", "quote-1").
					Return(&coreentity.InternalQuotation{ID: "quote-1", Status: coreentity.QuotationStatusRejected}, nil)
			},
			wantError: true,
		},
		{
			name: "converts approved quotation into order and invoice",
			input: coreentity.InternalQuotationActionInput{
				UserCtx: common.UserContext{TenantID: "tenant-1"},
				ID:      "quote-1",
				Action:  coreentity.ActionConvertQuotation,
			},
			setup: func(quoteRepo *dbMocks.InternalQuotationRepository, orderRepo *dbMocks.InternalOrderRepository, invoiceRepo *dbMocks.InternalInvoiceRepository, tx *dbMocks.Transactor) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				quoteRepo.EXPECT().
					GetQuotation(mock.Anything, "tenant-1", "quote-1").
					Return(&coreentity.InternalQuotation{
						ID:                       "quote-1",
						TenantID:                 "tenant-1",
						Status:                   coreentity.QuotationStatusApproved,
						CurrencyCode:             "IDR",
						TotalAmount:              decimal.NewFromInt(100),
						InternalProductID:        "product-1",
						InternalProductPricingID: "pricing-1",
					}, nil)
				orderRepo.EXPECT().
					CreateOrder(mock.Anything, mock.MatchedBy(func(data coreentity.InternalOrder) bool {
						return data.SourceType == coreentity.OrderSourceTypeQuotation &&
							data.SourceReferenceID != nil &&
							*data.SourceReferenceID == "quote-1" &&
							data.Status == coreentity.OrderStatusPendingPayment
					})).
					Return(&coreentity.InternalOrder{ID: "order-1", TotalAmount: decimal.NewFromInt(100), CurrencyCode: "IDR"}, nil)
				invoiceRepo.EXPECT().
					CreateInvoice(mock.Anything, mock.MatchedBy(func(data coreentity.InternalInvoice) bool {
						return data.InternalOrderID == "order-1" &&
							data.Status == coreentity.InvoiceStatusOpen &&
							data.AmountOutstanding.Equal(decimal.NewFromInt(100))
					})).
					Return(&coreentity.InternalInvoice{ID: "invoice-1"}, nil)
				quoteRepo.EXPECT().
					MarkQuotationConverted(mock.Anything, "tenant-1", "quote-1", "order-1").
					Return(&coreentity.InternalQuotation{ID: "quote-1", Status: coreentity.QuotationStatusConverted}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quoteRepo := dbMocks.NewInternalQuotationRepository(t)
			orderRepo := dbMocks.NewInternalOrderRepository(t)
			invoiceRepo := dbMocks.NewInternalInvoiceRepository(t)
			tx := dbMocks.NewTransactor(t)
			if tt.setup != nil {
				tt.setup(quoteRepo, orderRepo, invoiceRepo, tx)
			}
			core := NewInternalQuotationCore(Config{
				QuotationRepo: quoteRepo,
				OrderRepo:     orderRepo,
				InvoiceRepo:   invoiceRepo,
				Tx:            tx,
			})

			got, err := core.ExecuteQuotationAction(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got.Quotation)
		})
	}
}
