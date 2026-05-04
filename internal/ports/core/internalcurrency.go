package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalCurrencyCore interface {
	GetInternalCurrencies(
		ctx context.Context,
		filter coreentity.InternalCurrencyListFilter,
	) ([]coreentity.InternalCurrency, int, error)
	GetInternalCurrency(
		ctx context.Context,
		filter coreentity.InternalCurrencyFilter,
	) (*coreentity.InternalCurrency, error)
	CreateInternalCurrency(
		ctx context.Context,
		data coreentity.InternalCurrency,
	) (*coreentity.InternalCurrency, error)
	UpdateInternalCurrency(
		ctx context.Context,
		data coreentity.InternalCurrency,
	) error
	DeleteInternalCurrency(
		ctx context.Context,
		filter coreentity.InternalCurrencyDeleteFilter,
	) error

	GetProviderCurrencies(
		ctx context.Context,
		filter coreentity.InternalPaymentProviderCurrencyListFilter,
	) ([]coreentity.InternalPaymentProviderCurrency, int, error)
	UpsertProviderCurrency(
		ctx context.Context,
		data coreentity.InternalPaymentProviderCurrency,
	) error
	DeleteProviderCurrency(
		ctx context.Context,
		filter coreentity.InternalPaymentProviderCurrencyFilter,
	) error
}
