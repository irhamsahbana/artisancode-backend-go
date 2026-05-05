package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
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
