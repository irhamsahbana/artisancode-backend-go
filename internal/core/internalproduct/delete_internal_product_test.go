package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"
	"codebase-app/pkg/errmsg"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteInternalProduct(t *testing.T) {
	operatorCtx := context.WithValue(
		context.Background(),
		common.UserContextKeyClaims,
		common.UserContext{Roles: []string{coreentity.InternalUserRoleOperator}},
	)
	unauthorizedCtx := context.WithValue(
		context.Background(),
		common.UserContextKeyClaims,
		common.UserContext{Roles: []string{"viewer"}},
	)

	filter := coreentity.InternalProductDeleteFilter{ID: "product-1"}

	tests := []struct {
		name      string
		ctx       context.Context
		input     coreentity.InternalProductDeleteFilter
		setup     func(repo *dbMocks.InternalProductRepository)
		wantError bool
		assertErr func(t *testing.T, err error)
	}{
		{
			name:  "deletes product for authorized user",
			ctx:   operatorCtx,
			input: filter,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					DeleteInternalProduct(mock.Anything, coreentity.InternalProductDeleteFilter{ID: "product-1"}).
					Return(nil)
			},
		},
		{
			name:      "rejects unauthorized user before repository call",
			ctx:       unauthorizedCtx,
			input:     filter,
			wantError: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				customErr, ok := err.(*errmsg.CustomError)
				require.True(t, ok)
				require.Equal(t, errmsg.MessageYouAreNotAuthorizedToManageInternalProducts, customErr.Msg)
			},
		},
		{
			name:  "returns repository error",
			ctx:   operatorCtx,
			input: filter,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					DeleteInternalProduct(mock.Anything, coreentity.InternalProductDeleteFilter{ID: "product-1"}).
					Return(errmsg.NewCustomErrors(500).SetMessage("database error"))
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

			err := core.DeleteInternalProduct(tt.ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				if tt.assertErr != nil {
					tt.assertErr(t, err)
				}
				return
			}
			require.NoError(t, err)
		})
	}
}
