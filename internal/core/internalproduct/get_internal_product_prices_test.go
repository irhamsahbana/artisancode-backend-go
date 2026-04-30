package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetInternalProductPrices(t *testing.T) {
	repo := dbMocks.NewInternalProductRepository(t)
	repo.EXPECT().
		GetInternalProductPrices(
			mock.Anything,
			coreentity.InternalProductPriceListFilter{
				InternalProductPricingID: "pricing-1",
				CurrencyCode:             "IDR",
			},
		).
		Return([]coreentity.InternalProductPrice{{ID: "price-1"}}, 1, nil)
	core := NewInternalProductCore(Config{Repo: repo})

	got, total, err := core.GetInternalProductPrices(context.Background(), coreentity.InternalProductPriceListFilter{
		InternalProductPricingID: "pricing-1",
		CurrencyCode:             " idr ",
	})

	require.NoError(t, err)
	require.Equal(t, []coreentity.InternalProductPrice{{ID: "price-1"}}, got)
	require.Equal(t, 1, total)
}
