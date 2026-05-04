package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/coreentity"
)

func mapInternalCurrencyDAO(
	code string,
	name string,
	symbol string,
	decimalPlaces int,
	isActive bool,
	isDefault bool,
	sortOrder int,
	metadataRaw json.RawMessage,
	createdAt string,
	updatedAt *string,
) (coreentity.InternalCurrency, error) {
	metadata, err := parseMetadata(context.Background(), metadataRaw)
	if err != nil {
		return coreentity.InternalCurrency{}, err
	}

	item := coreentity.InternalCurrency{
		Code:          code,
		Name:          name,
		Symbol:        symbol,
		DecimalPlaces: decimalPlaces,
		IsActive:      isActive,
		IsDefault:     isDefault,
		SortOrder:     sortOrder,
		Metadata:      metadata,
		CreatedAt:     createdAt,
	}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}

	return item, nil
}

func mapInternalPaymentProviderCurrencyDAO(
	provider string,
	currencyCode string,
	isActive bool,
	minAmount *string,
	maxAmount *string,
	metadataRaw json.RawMessage,
	createdAt string,
	updatedAt *string,
) (coreentity.InternalPaymentProviderCurrency, error) {
	metadata, err := parseMetadata(context.Background(), metadataRaw)
	if err != nil {
		return coreentity.InternalPaymentProviderCurrency{}, err
	}

	item := coreentity.InternalPaymentProviderCurrency{
		Provider:     provider,
		CurrencyCode: currencyCode,
		IsActive:     isActive,
		MinAmount:    minAmount,
		MaxAmount:    maxAmount,
		Metadata:     metadata,
		CreatedAt:    createdAt,
	}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}

	return item, nil
}
