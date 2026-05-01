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

func TestUpdateInternalProductPrice(t *testing.T) {
	ctx := context.Background()
	userCtx := common.UserContext{Roles: []string{coreentity.InternalUserRoleOperator}}
	endedAt := " 2026-02-01T00:00:00Z "
	input := coreentity.InternalProductPrice{
		ID:                       "price-1",
		UserCtx:                  userCtx,
		InternalProductPricingID: "pricing-1",
		CurrencyCode:             " idr ",
		Amount:                   decimal.NewFromInt(175000),
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
			name:  "normalizes price and updates when authorized and period does not overlap",
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
							ExcludeID:                "price-1",
						},
					).
					Return(false, nil)
				repo.EXPECT().
					UpdateInternalProductPrice(
						mock.Anything,
						mock.MatchedBy(func(data coreentity.InternalProductPrice) bool {
							return data.ID == "price-1" &&
								data.CurrencyCode == "IDR" &&
								data.StartedAt == "2026-01-01T00:00:00Z" &&
								data.EndedAt != nil &&
								*data.EndedAt == "2026-02-01T00:00:00Z" &&
								data.Metadata != nil
						}),
					).
					Return(nil)
			},
		},
		{
			name: "rejects unauthorized user before repository call",
			input: coreentity.InternalProductPrice{
				ID:                       "price-1",
				UserCtx:                  common.UserContext{Roles: []string{"viewer"}},
				InternalProductPricingID: "pricing-1",
				CurrencyCode:             "IDR",
				Amount:                   decimal.NewFromInt(175000),
				StartedAt:                "2026-01-01T00:00:00Z",
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
			name:  "returns overlap lookup repository error",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProductPricing(mock.Anything, coreentity.InternalProductPricingFilter{ID: "pricing-1"}).
					Return(&coreentity.InternalProductPricing{ID: "pricing-1"}, nil)
				repo.EXPECT().
					ExistsOverlappingInternalProductPrice(mock.Anything, mock.Anything).
					Return(false, errors.New("repo failed"))
			},
			wantError: true,
		},
		{
			name:  "returns update repository error",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProductPricing(mock.Anything, coreentity.InternalProductPricingFilter{ID: "pricing-1"}).
					Return(&coreentity.InternalProductPricing{ID: "pricing-1"}, nil)
				repo.EXPECT().
					ExistsOverlappingInternalProductPrice(mock.Anything, mock.Anything).
					Return(false, nil)
				repo.EXPECT().
					UpdateInternalProductPrice(mock.Anything, mock.Anything).
					Return(errors.New("repo failed"))
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

			err := core.UpdateInternalProductPrice(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
