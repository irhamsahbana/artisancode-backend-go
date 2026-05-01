package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/pkg/errmsg"

	"github.com/stretchr/testify/require"
)

func TestNormalizeAndValidateProduct(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		input       coreentity.InternalProduct
		assert      func(t *testing.T, got coreentity.InternalProduct)
		wantCode    int
		wantMessage string
	}{
		{
			name: "normalizes valid product",
			input: coreentity.InternalProduct{
				Code:        " product_1 ",
				Name:        " Product 1 ",
				Description: " Description ",
				Status:      " ACTIVE ",
			},
			assert: func(t *testing.T, got coreentity.InternalProduct) {
				require.Equal(t, "PRODUCT_1", got.Code)
				require.Equal(t, "Product 1", got.Name)
				require.Equal(t, "Description", got.Description)
				require.Equal(t, coreentity.InternalProductStatusActive, got.Status)
				require.NotNil(t, got.Metadata)
				require.Empty(t, got.Metadata)
			},
		},
		{
			name: "rejects invalid code",
			input: coreentity.InternalProduct{
				Code:   "bad code",
				Status: coreentity.InternalProductStatusActive,
			},
			wantCode:    400,
			wantMessage: errmsg.MessageInternalProductCodeFormatIsInvalid,
		},
		{
			name: "rejects invalid status",
			input: coreentity.InternalProduct{
				Code:   "PRODUCT_1",
				Status: "published",
			},
			wantCode:    400,
			wantMessage: errmsg.MessageInternalProductStatusIsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := NewInternalProductCore(Config{})
			input := tt.input

			err := core.normalizeAndValidateProduct(ctx, &input)

			if tt.wantMessage == "" {
				require.NoError(t, err)
				tt.assert(t, input)
				return
			}

			requireHelperError(t, err, tt.wantCode, tt.wantMessage)
		})
	}
}
