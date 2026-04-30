package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	infraConfig "codebase-app/internal/infrastructure/config"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	infraConfig.Envs = &infraConfig.Config{}
	infraConfig.Envs.App.Name = "internalorder-core-test"
}

func TestCreateOrder(t *testing.T) {
	ctx := context.Background()
	companyID := "company-1"
	input := coreentity.CreateInternalOrderInput{
		UserCtx: common.UserContext{
			TenantID:  "tenant-1",
			CompanyID: &companyID,
		},
		InternalProductID:        "product-1",
		InternalProductPricingID: "pricing-1",
		CurrencyCode:             " idr ",
	}

	tests := []struct {
		name      string
		input     coreentity.CreateInternalOrderInput
		setup     func(orderRepo *dbMocks.InternalOrderRepository, invoiceRepo *dbMocks.InternalInvoiceRepository, tx *dbMocks.Transactor)
		wantError bool
	}{
		{
			name:  "creates order and invoice in transaction",
			input: input,
			setup: func(orderRepo *dbMocks.InternalOrderRepository, invoiceRepo *dbMocks.InternalInvoiceRepository, tx *dbMocks.Transactor) {
				orderRepo.EXPECT().
					GetPricingSnapshot(mock.Anything, "product-1", "pricing-1", "IDR").
					Return(map[string]any{"name": "Plan"}, "150000", nil)
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				orderRepo.EXPECT().
					CreateOrder(
						mock.Anything,
						mock.MatchedBy(func(data coreentity.InternalOrder) bool {
							return data.TenantID == "tenant-1" &&
								data.CompanyID != nil &&
								*data.CompanyID == "company-1" &&
								data.CurrencyCode == "IDR" &&
								data.TotalAmount.Equal(decimal.NewFromInt(150000)) &&
								data.Status == coreentity.OrderStatusPendingPayment
						}),
					).
					Return(&coreentity.InternalOrder{ID: "order-1", TotalAmount: decimal.NewFromInt(150000)}, nil)
				invoiceRepo.EXPECT().
					CreateInvoice(
						mock.Anything,
						mock.MatchedBy(func(data coreentity.InternalInvoice) bool {
							return data.InternalOrderID == "order-1" &&
								data.Status == coreentity.InvoiceStatusOpen &&
								data.Amount.Equal(decimal.NewFromInt(150000)) &&
								data.AmountOutstanding.Equal(decimal.NewFromInt(150000))
						}),
					).
					Return(&coreentity.InternalInvoice{ID: "invoice-1"}, nil)
			},
		},
		{
			name: "uses input tenant fallback",
			input: coreentity.CreateInternalOrderInput{
				UserCtx:                  common.UserContext{},
				TenantID:                 " tenant-fallback ",
				InternalProductID:        "product-1",
				InternalProductPricingID: "pricing-1",
				CurrencyCode:             "idr",
			},
			setup: func(orderRepo *dbMocks.InternalOrderRepository, invoiceRepo *dbMocks.InternalInvoiceRepository, tx *dbMocks.Transactor) {
				orderRepo.EXPECT().
					GetPricingSnapshot(mock.Anything, "product-1", "pricing-1", "IDR").
					Return(map[string]any{}, "1", nil)
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				orderRepo.EXPECT().
					CreateOrder(mock.Anything, mock.MatchedBy(func(data coreentity.InternalOrder) bool {
						return data.TenantID == "tenant-fallback"
					})).
					Return(&coreentity.InternalOrder{ID: "order-1", TotalAmount: decimal.NewFromInt(1)}, nil)
				invoiceRepo.EXPECT().
					CreateInvoice(mock.Anything, mock.Anything).
					Return(&coreentity.InternalInvoice{ID: "invoice-1"}, nil)
			},
		},
		{
			name: "rejects missing tenant before repository call",
			input: coreentity.CreateInternalOrderInput{
				InternalProductID:        "product-1",
				InternalProductPricingID: "pricing-1",
				CurrencyCode:             "IDR",
			},
			wantError: true,
		},
		{
			name:  "returns pricing lookup error",
			input: input,
			setup: func(orderRepo *dbMocks.InternalOrderRepository, invoiceRepo *dbMocks.InternalInvoiceRepository, tx *dbMocks.Transactor) {
				orderRepo.EXPECT().
					GetPricingSnapshot(mock.Anything, "product-1", "pricing-1", "IDR").
					Return(nil, "", errors.New("repo failed"))
			},
			wantError: true,
		},
		{
			name:  "returns invalid amount error",
			input: input,
			setup: func(orderRepo *dbMocks.InternalOrderRepository, invoiceRepo *dbMocks.InternalInvoiceRepository, tx *dbMocks.Transactor) {
				orderRepo.EXPECT().
					GetPricingSnapshot(mock.Anything, "product-1", "pricing-1", "IDR").
					Return(map[string]any{}, "not-a-number", nil)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderRepo := dbMocks.NewInternalOrderRepository(t)
			invoiceRepo := dbMocks.NewInternalInvoiceRepository(t)
			tx := dbMocks.NewTransactor(t)
			if tt.setup != nil {
				tt.setup(orderRepo, invoiceRepo, tx)
			}
			core := NewInternalOrderCore(Config{OrderRepo: orderRepo, InvoiceRepo: invoiceRepo, Tx: tx})

			got, err := core.CreateOrder(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got.Order)
			require.NotNil(t, got.Invoice)
		})
	}
}
