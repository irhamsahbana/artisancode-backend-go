package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/pkg/errmsg"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestNormalizeAndValidatePrice(t *testing.T) {
	ctx := context.Background()
	endedAt := " 2026-02-01T00:00:00Z "
	blankEndedAt := " "

	tests := []struct {
		name        string
		input       coreentity.InternalProductPrice
		assert      func(t *testing.T, got coreentity.InternalProductPrice)
		wantCode    int
		wantMessage string
	}{
		{
			name: "normalizes valid price",
			input: coreentity.InternalProductPrice{
				CurrencyCode: " idr ",
				Amount:       decimal.NewFromInt(150000),
				StartedAt:    " 2026-01-01T00:00:00Z ",
				EndedAt:      &endedAt,
			},
			assert: func(t *testing.T, got coreentity.InternalProductPrice) {
				require.Equal(t, "IDR", got.CurrencyCode)
				require.Equal(t, "2026-01-01T00:00:00Z", got.StartedAt)
				require.NotNil(t, got.EndedAt)
				require.Equal(t, "2026-02-01T00:00:00Z", *got.EndedAt)
				require.NotNil(t, got.Metadata)
				require.Empty(t, got.Metadata)
			},
		},
		{
			name: "normalizes blank ended at to nil",
			input: coreentity.InternalProductPrice{
				CurrencyCode: "USD",
				Amount:       decimal.NewFromInt(1),
				StartedAt:    "2026-01-01T00:00:00Z",
				EndedAt:      &blankEndedAt,
			},
			assert: func(t *testing.T, got coreentity.InternalProductPrice) {
				require.Nil(t, got.EndedAt)
				require.NotNil(t, got.Metadata)
			},
		},
		{
			name: "rejects invalid currency",
			input: coreentity.InternalProductPrice{
				CurrencyCode: "idr1",
				Amount:       decimal.NewFromInt(1),
				StartedAt:    "2026-01-01T00:00:00Z",
			},
			wantCode:    400,
			wantMessage: errmsg.MessageCurrencyCodeFormatIsInvalid,
		},
		{
			name: "rejects non positive amount",
			input: coreentity.InternalProductPrice{
				CurrencyCode: "IDR",
				Amount:       decimal.Zero,
				StartedAt:    "2026-01-01T00:00:00Z",
			},
			wantCode:    400,
			wantMessage: errmsg.MessageAmountMustBePositive,
		},
		{
			name: "rejects invalid started at",
			input: coreentity.InternalProductPrice{
				CurrencyCode: "IDR",
				Amount:       decimal.NewFromInt(1),
				StartedAt:    "2026-01-01",
			},
			wantCode:    400,
			wantMessage: errmsg.MessageStartedAtMustUseRfc3339Format,
		},
		{
			name: "rejects invalid ended at",
			input: coreentity.InternalProductPrice{
				CurrencyCode: "IDR",
				Amount:       decimal.NewFromInt(1),
				StartedAt:    "2026-01-01T00:00:00Z",
				EndedAt:      stringPtr("2026-01-02"),
			},
			wantCode:    400,
			wantMessage: errmsg.MessageEndedAtMustUseRfc3339Format,
		},
		{
			name: "rejects ended at not later than started at",
			input: coreentity.InternalProductPrice{
				CurrencyCode: "IDR",
				Amount:       decimal.NewFromInt(1),
				StartedAt:    "2026-01-01T00:00:00Z",
				EndedAt:      stringPtr("2026-01-01T00:00:00Z"),
			},
			wantCode:    400,
			wantMessage: errmsg.MessageEndedAtMustBeLaterThanStartedAt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := NewInternalProductCore(Config{})
			input := tt.input

			err := core.normalizeAndValidatePrice(ctx, &input)

			if tt.wantMessage == "" {
				require.NoError(t, err)
				tt.assert(t, input)
				return
			}

			requireHelperError(t, err, tt.wantCode, tt.wantMessage)
		})
	}
}

func stringPtr(value string) *string {
	return &value
}
