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

func TestNormalizeAndValidatePriceDecimalPlaces(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		input       coreentity.InternalProductPrice
		wantCode    int
		wantMessage string
	}{
		{
			name: "accepts zero decimal places for IDR",
			input: coreentity.InternalProductPrice{
				CurrencyCode: "IDR",
				Amount:       decimal.NewFromInt(150000),
				StartedAt:    "2026-01-01T00:00:00Z",
			},
		},
		{
			name: "rejects decimals for zero decimal currency",
			input: coreentity.InternalProductPrice{
				CurrencyCode: "IDR",
				Amount:       decimal.RequireFromString("150000.50"),
				StartedAt:    "2026-01-01T00:00:00Z",
			},
			wantCode:    400,
			wantMessage: errmsg.MessageCurrencyAmountPrecisionIsInvalid,
		},
		{
			name: "accepts two decimal places for USD",
			input: coreentity.InternalProductPrice{
				CurrencyCode: "USD",
				Amount:       decimal.RequireFromString("19.99"),
				StartedAt:    "2026-01-01T00:00:00Z",
			},
		},
		{
			name: "rejects three decimal places for two decimal currency",
			input: coreentity.InternalProductPrice{
				CurrencyCode: "USD",
				Amount:       decimal.RequireFromString("19.999"),
				StartedAt:    "2026-01-01T00:00:00Z",
			},
			wantCode:    400,
			wantMessage: errmsg.MessageCurrencyAmountPrecisionIsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := NewInternalProductCore(Config{CurrencyRepo: currencyRepoStub{}})
			input := tt.input

			err := core.normalizeAndValidatePrice(ctx, &input)

			if tt.wantMessage == "" {
				require.NoError(t, err)
				return
			}

			requireHelperError(t, err, tt.wantCode, tt.wantMessage)
		})
	}
}

func stringPtr(value string) *string {
	return &value
}

type currencyRepoStub struct{}

func (currencyRepoStub) GetInternalCurrencies(
	ctx context.Context,
	filter coreentity.InternalCurrencyListFilter,
) ([]coreentity.InternalCurrency, int, error) {
	return nil, 0, nil
}

func (currencyRepoStub) GetInternalCurrency(
	ctx context.Context,
	filter coreentity.InternalCurrencyFilter,
) (*coreentity.InternalCurrency, error) {
	decimalPlaces := 0
	if filter.Code == "USD" || filter.Code == "EUR" || filter.Code == "SGD" {
		decimalPlaces = 2
	}
	return &coreentity.InternalCurrency{
		Code:          filter.Code,
		Symbol:        filter.Code,
		DecimalPlaces: decimalPlaces,
		IsActive:      true,
		IsDefault:     filter.Code == "IDR",
	}, nil
}

func (currencyRepoStub) CreateInternalCurrency(
	ctx context.Context,
	data coreentity.InternalCurrency,
) (*coreentity.InternalCurrency, error) {
	return nil, nil
}

func (currencyRepoStub) UpdateInternalCurrency(
	ctx context.Context,
	data coreentity.InternalCurrency,
) error {
	return nil
}

func (currencyRepoStub) DeleteInternalCurrency(
	ctx context.Context,
	filter coreentity.InternalCurrencyDeleteFilter,
) error {
	return nil
}

func (currencyRepoStub) IsCurrencyActive(ctx context.Context, code string) (bool, error) {
	return true, nil
}

func (currencyRepoStub) GetDefaultCurrency(ctx context.Context) (*coreentity.InternalCurrency, error) {
	return &coreentity.InternalCurrency{Code: "IDR", IsActive: true, IsDefault: true}, nil
}
