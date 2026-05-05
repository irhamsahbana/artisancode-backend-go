package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalCurrencyRepository interface {
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
	IsCurrencyActive(ctx context.Context, code string) (bool, error)
	GetDefaultCurrency(ctx context.Context) (*coreentity.InternalCurrency, error)
}
