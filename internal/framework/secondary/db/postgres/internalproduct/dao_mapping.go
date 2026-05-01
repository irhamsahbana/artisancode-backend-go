package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/coreentity"

	"github.com/shopspring/decimal"
)

func mapInternalProductDAO(
	id string,
	code string,
	name string,
	description string,
	status string,
	metadataRaw json.RawMessage,
	createdAt string,
	updatedAt *string,
) (coreentity.InternalProduct, error) {
	metadata, err := parseMetadata(context.Background(), metadataRaw)
	if err != nil {
		return coreentity.InternalProduct{}, err
	}

	item := coreentity.InternalProduct{
		ID:          id,
		Code:        code,
		Name:        name,
		Description: description,
		Status:      status,
		Metadata:    metadata,
		CreatedAt:   createdAt,
	}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}

	return item, nil
}

func mapInternalProductPricingDAO(
	id string,
	internalProductID string,
	code string,
	name string,
	description string,
	status string,
	metadataRaw json.RawMessage,
	createdAt string,
	updatedAt *string,
) (coreentity.InternalProductPricing, error) {
	metadata, err := parseMetadata(context.Background(), metadataRaw)
	if err != nil {
		return coreentity.InternalProductPricing{}, err
	}

	item := coreentity.InternalProductPricing{
		ID:                id,
		InternalProductID: internalProductID,
		Code:              code,
		Name:              name,
		Description:       description,
		Status:            status,
		Metadata:          metadata,
		CreatedAt:         createdAt,
	}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}

	return item, nil
}

func mapInternalProductPriceDAO(
	id string,
	pricingID string,
	currencyCode string,
	amount decimal.Decimal,
	startedAt string,
	endedAt *string,
	metadataRaw json.RawMessage,
	createdAt string,
	updatedAt *string,
) (coreentity.InternalProductPrice, error) {
	metadata, err := parseMetadata(context.Background(), metadataRaw)
	if err != nil {
		return coreentity.InternalProductPrice{}, err
	}

	item := coreentity.InternalProductPrice{
		ID:                       id,
		InternalProductPricingID: pricingID,
		CurrencyCode:             currencyCode,
		Amount:                   amount,
		StartedAt:                startedAt,
		EndedAt:                  endedAt,
		Metadata:                 metadata,
		CreatedAt:                createdAt,
	}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}

	return item, nil
}
