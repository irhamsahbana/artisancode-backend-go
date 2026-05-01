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

func TestGetInternalProducts(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.InternalProductListFilter{
		Q:        "  Product Query  ",
		Status:   " ACTIVE ",
		Page:     2,
		Paginate: 10,
	}

	tests := []struct {
		name      string
		setup     func(repo *dbMocks.InternalProductRepository)
		want      []coreentity.InternalProduct
		wantTotal int
		wantError bool
	}{
		{
			name: "passes filter through and returns products",
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProducts(mock.Anything, filter).
					Return([]coreentity.InternalProduct{
						{ID: "product-1", Code: "PRODUCT_1"},
					}, 1, nil)
			},
			want:      []coreentity.InternalProduct{{ID: "product-1", Code: "PRODUCT_1"}},
			wantTotal: 1,
		},
		{
			name: "returns repository error",
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProducts(mock.Anything, filter).
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

			got, total, err := core.GetInternalProducts(ctx, filter)

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

func TestGetInternalProduct(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.InternalProductFilter{ID: " product-1 "}

	tests := []struct {
		name      string
		setup     func(repo *dbMocks.InternalProductRepository)
		want      *coreentity.InternalProduct
		wantError bool
	}{
		{
			name: "passes filter through and returns product",
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProduct(mock.Anything, filter).
					Return(&coreentity.InternalProduct{ID: "product-1", Code: "PRODUCT_1"}, nil)
			},
			want: &coreentity.InternalProduct{ID: "product-1", Code: "PRODUCT_1"},
		},
		{
			name: "returns repository error",
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					GetInternalProduct(mock.Anything, filter).
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

			got, err := core.GetInternalProduct(ctx, filter)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
