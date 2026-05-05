package core

import (
	"context"
	"strings"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *internalProductCore) normalizeAndValidatePrice(
	ctx context.Context,
	data *coreentity.InternalProductPrice,
) error {
	data.CurrencyCode = strings.ToUpper(strings.TrimSpace(data.CurrencyCode))
	data.StartedAt = strings.TrimSpace(data.StartedAt)
	if data.EndedAt != nil {
		trimmed := strings.TrimSpace(*data.EndedAt)
		data.EndedAt = &trimmed
		if trimmed == "" {
			data.EndedAt = nil
		}
	}
	if data.Metadata == nil {
		data.Metadata = map[string]any{}
	}

	if !internalCurrencyPattern.MatchString(data.CurrencyCode) {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageCurrencyCodeFormatIsInvalid)
	}

	var currencyDecimalPlaces int
	if c.currencyRepo != nil {
		currency, err := c.currencyRepo.GetInternalCurrency(ctx, coreentity.InternalCurrencyFilter{
			Code: data.CurrencyCode,
		})
		if err != nil {
			return err
		}
		if !currency.IsActive {
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageCurrencyIsNotActive)
		}
		currencyDecimalPlaces = currency.DecimalPlaces
	}

	if !data.Amount.IsPositive() {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageAmountMustBePositive)
	}

	if currencyDecimalPlaces > 0 {
		truncated := data.Amount.Truncate(int32(currencyDecimalPlaces))
		if !data.Amount.Equal(truncated) {
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageCurrencyAmountPrecisionIsInvalid)
		}
	} else {
		if !data.Amount.Equal(data.Amount.Truncate(0)) {
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageCurrencyAmountPrecisionIsInvalid)
		}
	}

	startTime, err := time.Parse(time.RFC3339, data.StartedAt)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, data.StartedAt).Msg("Invalid started_at")
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageStartedAtMustUseRfc3339Format)
	}

	if data.EndedAt != nil {
		endTime, err := time.Parse(time.RFC3339, *data.EndedAt)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, data.EndedAt).Msg("Invalid ended_at")
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageEndedAtMustUseRfc3339Format)
		}
		if !endTime.After(startTime) {
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageEndedAtMustBeLaterThanStartedAt)
		}
	}

	return nil
}
