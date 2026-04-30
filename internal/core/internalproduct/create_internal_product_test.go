package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	infraConfig "codebase-app/internal/infrastructure/config"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	infraConfig.Envs = &infraConfig.Config{}
	infraConfig.Envs.App.Name = "internalproduct-core-test"
}

func TestCreateInternalProduct(t *testing.T) {
	ctx := context.Background()
	userCtx := common.UserContext{Roles: []string{coreentity.InternalUserRoleOperator}}
	input := coreentity.InternalProduct{
		UserCtx:     userCtx,
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
		want      *coreentity.InternalProduct
		wantError bool
	}{
		{
			name:  "normalizes data and creates product",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					ExistsInternalProductByCode(mock.Anything, "PRODUCT_1", "").
					Return(false, nil)
				repo.EXPECT().
					CreateInternalProduct(
						mock.Anything,
						mock.MatchedBy(func(data coreentity.InternalProduct) bool {
							return data.Code == "PRODUCT_1" &&
								data.Name == "Product 1" &&
								data.Description == "Description" &&
								data.Status == coreentity.InternalProductStatusActive &&
								data.Metadata != nil
						}),
					).
					Return(&coreentity.InternalProduct{ID: "product-1", Code: "PRODUCT_1"}, nil)
			},
			want: &coreentity.InternalProduct{ID: "product-1", Code: "PRODUCT_1"},
		},
		{
			name: "rejects unauthorized user before repository call",
			input: coreentity.InternalProduct{
				UserCtx: common.UserContext{Roles: []string{"viewer"}},
				Code:    "PRODUCT_1",
				Status:  coreentity.InternalProductStatusActive,
			},
			wantError: true,
		},
		{
			name: "rejects invalid code before repository call",
			input: coreentity.InternalProduct{
				UserCtx: userCtx,
				Code:    "bad code",
				Status:  coreentity.InternalProductStatusActive,
			},
			wantError: true,
		},
		{
			name:  "rejects duplicate product code",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					ExistsInternalProductByCode(mock.Anything, "PRODUCT_1", "").
					Return(true, nil)
			},
			wantError: true,
		},
		{
			name:  "returns repository error",
			input: input,
			setup: func(repo *dbMocks.InternalProductRepository) {
				repo.EXPECT().
					ExistsInternalProductByCode(mock.Anything, "PRODUCT_1", "").
					Return(false, errors.New("repo failed"))
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

			got, err := core.CreateInternalProduct(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
