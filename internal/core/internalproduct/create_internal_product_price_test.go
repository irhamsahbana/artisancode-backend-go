package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateInternalProductPrice(t *testing.T) {
	ctx := context.Background()
	userCtx := common.UserContext{Roles: []string{coreentity.InternalUserRoleSuperAdmin}}
	endedAt := " 2026-02-01T00:00:00Z "
	input := coreentity.InternalProductPrice{
		UserCtx:                  userCtx,
		InternalProductPricingID: "pricing-1",
		CurrencyCode:             " idr ",
		Amount:                   decimal.NewFromInt(150000),
		StartedAt:                " 2026-01-01T00:00:00Z ",
		EndedAt:                  &endedAt,
	}

	tests := []struct {
		name      string
		input     coreentity.InternalProductPrice
		setup     func(repo *dbMocks.InternalProductRepository)
		wantError bool
	}{
		{
			name:  "normalizes price and creates when period does not overlap",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProductPricing(mock.Anything, coreentity.InternalProductPricingFilter{ID: "pricing-1"}).
					Return(&coreentity.InternalProductPricing{ID: "pricing-1"}, nil)
				repo.EXPECT().
					ExistsOverlappingInternalProductPrice(
						mock.Anything,
						coreentity.InternalProductPriceOverlapFilter{
							InternalProductPricingID: "pricing-1",
							CurrencyCode:             "IDR",
							StartedAt:                "2026-01-01T00:00:00Z",
							EndedAt:                  ptr("2026-02-01T00:00:00Z"),
						},
					).
					Return(false, nil)
				repo.EXPECT().
					CreateInternalProductPrice(
						mock.Anything,
						mock.MatchedBy(func(data coreentity.InternalProductPrice) bool {
							return data.CurrencyCode == "IDR" &&
								data.StartedAt == "2026-01-01T00:00:00Z" &&
								data.EndedAt != nil &&
								*data.EndedAt == "2026-02-01T00:00:00Z" &&
								data.Metadata != nil
						}),
					).
					Return(&coreentity.InternalProductPrice{ID: "price-1"}, nil)
			},
		},
		{
			name: "rejects non-positive amount",
			input: coreentity.InternalProductPrice{
				UserCtx:                  userCtx,
				InternalProductPricingID: "pricing-1",
				CurrencyCode:             "IDR",
				Amount:                   decimal.Zero,
				StartedAt:                "2026-01-01T00:00:00Z",
			},
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProductPricing(mock.Anything, coreentity.InternalProductPricingFilter{ID: "pricing-1"}).
					Return(&coreentity.InternalProductPricing{ID: "pricing-1"}, nil)
			},
			wantError: true,
		},
		{
			name:  "rejects overlapping active price",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProductPricing(mock.Anything, coreentity.InternalProductPricingFilter{ID: "pricing-1"}).
					Return(&coreentity.InternalProductPricing{ID: "pricing-1"}, nil)
				repo.EXPECT().
					ExistsOverlappingInternalProductPrice(mock.Anything, mock.Anything).
					Return(true, nil)
			},
			wantError: true,
		},
		{
			name:  "returns pricing lookup error",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProductPricing(mock.Anything, coreentity.InternalProductPricingFilter{ID: "pricing-1"}).
					Return(nil, errors.New("repo failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewInternalProductRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}
			core := NewInternalProductCore(Config{Repo: repo})

			got, err := core.CreateInternalProductPrice(ctx, tt.input)

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

func ptr(value string) *string {
	return &value
}
