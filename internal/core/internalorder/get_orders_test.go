package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetOrders(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{TenantID: "tenant-1"})
	repo := dbMocks.NewInternalOrderRepository(t)
	repo.EXPECT().
		GetOrders(mock.Anything, coreentity.InternalCommerceListFilter{TenantID: "tenant-1"}).
		Return([]coreentity.InternalOrder{{ID: "order-1"}}, 1, nil)
	core := NewInternalOrderCore(Config{OrderRepo: repo})

	got, total, err := core.GetOrders(ctx, coreentity.InternalCommerceListFilter{})

	require.NoError(t, err)
	require.Equal(t, []coreentity.InternalOrder{{ID: "order-1"}}, got)
	require.Equal(t, 1, total)
}
