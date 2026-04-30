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
	infraConfig.Envs.App.Name = "internalquotation-core-test"
}

func TestCreateQuotation(t *testing.T) {
	ctx := context.Background()
	companyID := "company-1"
	input := coreentity.CreateInternalQuotationInput{
		UserCtx: common.UserContext{
			TenantID:  "tenant-1",
			CompanyID: &companyID,
		},
		InternalProductID:        "product-1",
		InternalProductPricingID: "pricing-1",
		CurrencyCode:             " idr ",
		SubtotalAmount:           decimal.NewFromInt(100),
		DiscountAmount:           decimal.NewFromInt(10),
		TaxAmount:                decimal.NewFromInt(1),
		TotalAmount:              decimal.NewFromInt(91),
		QuoteSnapshot:            map[string]any{"seat": 1},
	}

	tests := []struct {
		name      string
		input     coreentity.CreateInternalQuotationInput
		setup     func(orderRepo *dbMocks.InternalOrderRepository, quotationRepo *dbMocks.InternalQuotationRepository)
		wantError bool
	}{
		{
			name:  "creates draft quotation with normalized currency",
			input: input,
			setup: func(orderRepo *dbMocks.InternalOrderRepository, quotationRepo *dbMocks.InternalQuotationRepository) {
				orderRepo.EXPECT().
					GetPricingSnapshot(mock.Anything, "product-1", "pricing-1", "IDR").
					Return(map[string]any{"name": "Plan"}, "100", nil)
				quotationRepo.EXPECT().
					CreateQuotation(
						mock.Anything,
						mock.MatchedBy(func(data coreentity.InternalQuotation) bool {
							return data.TenantID == "tenant-1" &&
								data.CompanyID != nil &&
								*data.CompanyID == "company-1" &&
								data.Status == coreentity.QuotationStatusDraft &&
								data.CurrencyCode == "IDR" &&
								data.PricingSnapshot["name"] == "Plan"
						}),
					).
					Return(&coreentity.InternalQuotation{ID: "quote-1"}, nil)
			},
		},
		{
			name: "uses input tenant fallback",
			input: coreentity.CreateInternalQuotationInput{
				TenantID:                 " tenant-fallback ",
				InternalProductID:        "product-1",
				InternalProductPricingID: "pricing-1",
				CurrencyCode:             "idr",
			},
			setup: func(orderRepo *dbMocks.InternalOrderRepository, quotationRepo *dbMocks.InternalQuotationRepository) {
				orderRepo.EXPECT().
					GetPricingSnapshot(mock.Anything, "product-1", "pricing-1", "IDR").
					Return(map[string]any{}, "100", nil)
				quotationRepo.EXPECT().
					CreateQuotation(mock.Anything, mock.MatchedBy(func(data coreentity.InternalQuotation) bool {
						return data.TenantID == "tenant-fallback"
					})).
					Return(&coreentity.InternalQuotation{ID: "quote-1"}, nil)
			},
		},
		{
			name:      "rejects missing tenant",
			input:     coreentity.CreateInternalQuotationInput{CurrencyCode: "IDR"},
			wantError: true,
		},
		{
			name: "rejects negative amount",
			input: coreentity.CreateInternalQuotationInput{
				UserCtx:        common.UserContext{TenantID: "tenant-1"},
				CurrencyCode:   "IDR",
				DiscountAmount: decimal.NewFromInt(-1),
			},
			wantError: true,
		},
		{
			name:  "returns pricing lookup error",
			input: input,
			setup: func(orderRepo *dbMocks.InternalOrderRepository, quotationRepo *dbMocks.InternalQuotationRepository) {
				orderRepo.EXPECT().
					GetPricingSnapshot(mock.Anything, "product-1", "pricing-1", "IDR").
					Return(nil, "", errors.New("repo failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderRepo := dbMocks.NewInternalOrderRepository(t)
			quotationRepo := dbMocks.NewInternalQuotationRepository(t)
			if tt.setup != nil {
				tt.setup(orderRepo, quotationRepo)
			}
			core := NewInternalQuotationCore(Config{QuotationRepo: quotationRepo, OrderRepo: orderRepo})

			got, err := core.CreateQuotation(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
