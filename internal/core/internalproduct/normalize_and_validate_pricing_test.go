package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/pkg/errmsg"

	"github.com/stretchr/testify/require"
)

func TestNormalizeAndValidatePricing(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		input       coreentity.InternalProductPricing
		assert      func(t *testing.T, got coreentity.InternalProductPricing)
		wantCode    int
		wantMessage string
	}{
		{
			name: "normalizes valid pricing",
			input: coreentity.InternalProductPricing{
				Code:        " pricing_1 ",
				Name:        " Pricing 1 ",
				Description: " Description ",
				Status:      " DRAFT ",
			},
			assert: func(t *testing.T, got coreentity.InternalProductPricing) {
				require.Equal(t, "PRICING_1", got.Code)
				require.Equal(t, "Pricing 1", got.Name)
				require.Equal(t, "Description", got.Description)
				require.Equal(t, coreentity.InternalProductStatusDraft, got.Status)
				require.NotNil(t, got.Metadata)
				require.Empty(t, got.Metadata)
			},
		},
		{
			name: "rejects invalid code",
			input: coreentity.InternalProductPricing{
				Code:   "bad code",
				Status: coreentity.InternalProductStatusActive,
			},
			wantCode:    400,
			wantMessage: errmsg.MessageInternalProductPricingCodeFormatIsInvalid,
		},
		{
			name: "rejects invalid status",
			input: coreentity.InternalProductPricing{
				Code:   "PRICING_1",
				Status: "published",
			},
			wantCode:    400,
			wantMessage: errmsg.MessageInternalProductPricingStatusIsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := NewInternalProductCore(Config{})
			input := tt.input

			err := core.normalizeAndValidatePricing(ctx, &input)

			if tt.wantMessage == "" {
				require.NoError(t, err)
				tt.assert(t, input)
				return
			}

			requireHelperError(t, err, tt.wantCode, tt.wantMessage)
		})
	}
}
