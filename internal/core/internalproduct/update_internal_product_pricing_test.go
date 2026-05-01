package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateInternalProductPricing(t *testing.T) {
	ctx := context.Background()
	userCtx := common.UserContext{Roles: []string{coreentity.InternalUserRoleSuperAdmin}}
	input := coreentity.InternalProductPricing{
		ID:                "pricing-1",
		UserCtx:           userCtx,
		InternalProductID: "product-1",
		Code:              " pricing_1 ",
		Name:              " Pricing 1 ",
		Description:       " Description ",
		Status:            " INACTIVE ",
	}

	tests := []struct {
		name      string
		input     coreentity.InternalProductPricing
		setup     func(repo *dbMocks.InternalProductRepository)
		wantError bool
	}{
		{
			name:  "normalizes data and updates pricing",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProduct(mock.Anything, coreentity.InternalProductFilter{ID: "product-1"}).
					Return(&coreentity.InternalProduct{ID: "product-1"}, nil)
				repo.EXPECT().
					ExistsInternalProductPricingByCode(mock.Anything, "product-1", "PRICING_1", "pricing-1").
					Return(false, nil)
				repo.EXPECT().
					UpdateInternalProductPricing(
						mock.Anything,
						mock.MatchedBy(func(data coreentity.InternalProductPricing) bool {
							return data.ID == "pricing-1" &&
								data.Code == "PRICING_1" &&
								data.Name == "Pricing 1" &&
								data.Description == "Description" &&
								data.Status == coreentity.InternalProductStatusInactive &&
								data.Metadata != nil
						}),
					).
					Return(nil)
			},
		},
		{
			name: "rejects unauthorized user before repository call",
			input: coreentity.InternalProductPricing{
				ID:                "pricing-1",
				UserCtx:           common.UserContext{Roles: []string{"viewer"}},
				InternalProductID: "product-1",
				Code:              "PRICING_1",
				Status:            coreentity.InternalProductStatusActive,
			},
			wantError: true,
		},
		{
			name:  "rejects duplicate pricing code",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProduct(mock.Anything, coreentity.InternalProductFilter{ID: "product-1"}).
					Return(&coreentity.InternalProduct{ID: "product-1"}, nil)
				repo.EXPECT().
					ExistsInternalProductPricingByCode(mock.Anything, "product-1", "PRICING_1", "pricing-1").
					Return(true, nil)
			},
			wantError: true,
		},
		{
			name:  "returns parent product lookup error",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProduct(mock.Anything, coreentity.InternalProductFilter{ID: "product-1"}).
					Return(nil, errors.New("repo failed"))
			},
			wantError: true,
		},
		{
			name:  "returns duplicate check repository error",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProduct(mock.Anything, coreentity.InternalProductFilter{ID: "product-1"}).
					Return(&coreentity.InternalProduct{ID: "product-1"}, nil)
				repo.EXPECT().
					ExistsInternalProductPricingByCode(mock.Anything, "product-1", "PRICING_1", "pricing-1").
					Return(false, errors.New("repo failed"))
			},
			wantError: true,
		},
		{
			name:  "returns update repository error",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProduct(mock.Anything, coreentity.InternalProductFilter{ID: "product-1"}).
					Return(&coreentity.InternalProduct{ID: "product-1"}, nil)
				repo.EXPECT().
					ExistsInternalProductPricingByCode(mock.Anything, "product-1", "PRICING_1", "pricing-1").
					Return(false, nil)
				repo.EXPECT().
					UpdateInternalProductPricing(mock.Anything, mock.Anything).
					Return(errors.New("repo failed"))
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

			err := core.UpdateInternalProductPricing(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
