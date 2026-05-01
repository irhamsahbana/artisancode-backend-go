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

func TestDeleteInternalProductPricing(t *testing.T) {
	filter := coreentity.InternalProductPricingDeleteFilter{ID: "pricing-1"}

	tests := []struct {
		name      string
		ctx       context.Context
		setup     func(repo *dbMocks.InternalProductRepository)
		wantError bool
	}{
		{
			name: "deletes pricing when authorized",
			ctx: context.WithValue(
				context.Background(),
				common.UserContextKeyClaims,
				common.UserContext{Roles: []string{coreentity.InternalUserRoleOperator}},
			),
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					DeleteInternalProductPricing(mock.Anything, filter).
					Return(nil)
			},
		},
		{
			name: "rejects unauthorized user before repository call",
			ctx: context.WithValue(
				context.Background(),
				common.UserContextKeyClaims,
				common.UserContext{Roles: []string{"viewer"}},
			),
			wantError: true,
		},
		{
			name: "returns repository error",
			ctx: context.WithValue(
				context.Background(),
				common.UserContextKeyClaims,
				common.UserContext{Roles: []string{coreentity.InternalUserRoleSuperAdmin}},
			),
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					DeleteInternalProductPricing(mock.Anything, filter).
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

			err := core.DeleteInternalProductPricing(tt.ctx, filter)

			if tt.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
