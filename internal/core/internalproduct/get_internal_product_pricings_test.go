package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetInternalProductPricings(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.InternalProductPricingListFilter{
		InternalProductID: " product-1 ",
		Q:                 "  Pricing Query  ",
		Page:              3,
		Paginate:          5,
	}

	tests := []struct {
		name      string
		setup     func(repo *dbMocks.InternalProductRepository)
		want      []coreentity.InternalProductPricing
		wantTotal int
		wantError bool
	}{
		{
			name: "passes filter through and returns pricings",
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProductPricings(mock.Anything, filter).
					Return([]coreentity.InternalProductPricing{
						{ID: "pricing-1", InternalProductID: "product-1", Code: "BASIC"},
					}, 1, nil)
			},
			want: []coreentity.InternalProductPricing{
				{ID: "pricing-1", InternalProductID: "product-1", Code: "BASIC"},
			},
			wantTotal: 1,
		},
		{
			name: "returns repository error",
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProductPricings(mock.Anything, filter).
					Return(nil, 0, errors.New("repo failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewInternalProductRepository(t)
			tt.setup(repo)
			core := NewInternalProductCore(Config{Repo: repo})

			got, total, err := core.GetInternalProductPricings(ctx, filter)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantTotal, total)
		})
	}
}

func TestGetInternalProductPricing(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.InternalProductPricingFilter{ID: " pricing-1 "}

	tests := []struct {
		name      string
		setup     func(repo *dbMocks.InternalProductRepository)
		want      *coreentity.InternalProductPricing
		wantError bool
	}{
		{
			name: "passes filter through and returns pricing",
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProductPricing(mock.Anything, filter).
					Return(&coreentity.InternalProductPricing{
						ID:                "pricing-1",
						InternalProductID: "product-1",
						Code:              "BASIC",
					}, nil)
			},
			want: &coreentity.InternalProductPricing{
				ID:                "pricing-1",
				InternalProductID: "product-1",
				Code:              "BASIC",
			},
		},
		{
			name: "returns repository error",
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProductPricing(mock.Anything, filter).
					Return(nil, errors.New("repo failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewInternalProductRepository(t)
			tt.setup(repo)
			core := NewInternalProductCore(Config{Repo: repo})

			got, err := core.GetInternalProductPricing(ctx, filter)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
