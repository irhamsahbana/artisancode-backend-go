package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

func (c *internalCurrencyCore) normalizeAndValidateCurrency(
	ctx context.Context,
	data *coreentity.InternalCurrency,
) error {
	data.Code = strings.ToUpper(strings.TrimSpace(data.Code))
	data.Name = strings.TrimSpace(data.Name)
	data.Symbol = strings.TrimSpace(data.Symbol)
	if data.Metadata == nil {
		data.Metadata = map[string]any{}
	}

	if !currencyCodePattern.MatchString(data.Code) {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data.Code).Msg("Invalid internal currency code")
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageCurrencyCodeFormatIsInvalid)
	}
	if data.IsDefault {
		data.IsActive = true
	}
	if data.Name == "" {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageValidationRequired)
	}

	return nil
}

func (c *internalCurrencyCore) normalizeAndValidateProviderCurrency(
	ctx context.Context,
	data *coreentity.InternalPaymentProviderCurrency,
) error {
	data.Provider = strings.ToLower(strings.TrimSpace(data.Provider))
	data.CurrencyCode = strings.ToUpper(strings.TrimSpace(data.CurrencyCode))
	if data.Metadata == nil {
		data.Metadata = map[string]any{}
	}
	if !providerPattern.MatchString(data.Provider) {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data.Provider).Msg("Invalid provider format")
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageValidationSimpleFormat)
	}
	if !currencyCodePattern.MatchString(data.CurrencyCode) {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageCurrencyCodeFormatIsInvalid)
	}

	if data.MinAmount != nil {
		trimmed := strings.TrimSpace(*data.MinAmount)
		data.MinAmount = &trimmed
		if trimmed == "" {
			data.MinAmount = nil
		}
	}
	if data.MaxAmount != nil {
		trimmed := strings.TrimSpace(*data.MaxAmount)
		data.MaxAmount = &trimmed
		if trimmed == "" {
			data.MaxAmount = nil
		}
	}

	var minValue *decimal.Decimal
	if data.MinAmount != nil {
		value, err := decimal.NewFromString(*data.MinAmount)
		if err != nil || value.IsNegative() {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data.MinAmount).Msg("Invalid provider currency min amount")
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageValidationNumeric)
		}
		minValue = &value
	}

	if data.MaxAmount != nil {
		value, err := decimal.NewFromString(*data.MaxAmount)
		if err != nil || value.IsNegative() {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data.MaxAmount).Msg("Invalid provider currency max amount")
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageValidationNumeric)
		}
		if minValue != nil && value.LessThan(*minValue) {
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageValidationCompareGte)
		}
	}

	return nil
}
