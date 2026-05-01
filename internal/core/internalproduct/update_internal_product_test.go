package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"
	"codebase-app/pkg/errmsg"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateInternalProduct(t *testing.T) {
	ctx := context.Background()
	userCtx := common.UserContext{Roles: []string{coreentity.InternalUserRoleOperator}}
	input := coreentity.InternalProduct{
		UserCtx:     userCtx,
		ID:          "product-1",
		Code:        " product_1 ",
		Name:        " Product 1 ",
		Description: " Description ",
		Status:      " ACTIVE ",
		Metadata:    nil,
	}

	tests := []struct {
		name      string
		input     coreentity.InternalProduct
		setup     func(repo *dbMocks.InternalProductRepository)
		wantError bool
		assertErr func(t *testing.T, err error)
	}{
		{
			name:  "normalizes data and updates product",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					ExistsInternalProductByCode(mock.Anything, "PRODUCT_1", "product-1").
					Return(false, nil)
				repo.EXPECT().
					UpdateInternalProduct(
						mock.Anything,
						mock.MatchedBy(func(data coreentity.InternalProduct) bool {
							return data.ID == "product-1" &&
								data.Code == "PRODUCT_1" &&
								data.Name == "Product 1" &&
								data.Description == "Description" &&
								data.Status == coreentity.InternalProductStatusActive &&
								data.Metadata != nil
						}),
					).
					Return(nil)
			},
		},
		{
			name: "rejects unauthorized user before repository call",
			input: coreentity.InternalProduct{
				UserCtx: common.UserContext{Roles: []string{"viewer"}},
				ID:      "product-1",
				Code:    "PRODUCT_1",
				Status:  coreentity.InternalProductStatusActive,
			},
			wantError: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				customErr, ok := err.(*errmsg.CustomError)
				require.True(t, ok)
				require.Equal(t, errmsg.MessageYouAreNotAuthorizedToManageInternalProducts, customErr.Msg)
			},
		},
		{
			name:  "rejects duplicate product code",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					ExistsInternalProductByCode(mock.Anything, "PRODUCT_1", "product-1").
					Return(true, nil)
			},
			wantError: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				customErr, ok := err.(*errmsg.CustomError)
				require.True(t, ok)
				require.Equal(t, errmsg.MessageInternalProductCodeAlreadyExists, customErr.Msg)
			},
		},
		{
			name:  "returns repository error from duplicate check",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					ExistsInternalProductByCode(mock.Anything, "PRODUCT_1", "product-1").
					Return(false, errors.New("repo failed"))
			},
			wantError: true,
		},
		{
			name:  "returns repository error from update",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					ExistsInternalProductByCode(mock.Anything, "PRODUCT_1", "product-1").
					Return(false, nil)
				repo.EXPECT().
					UpdateInternalProduct(mock.Anything, mock.Anything).
					Return(errors.New("update failed"))
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

			err := core.UpdateInternalProduct(ctx, tt.input)

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
