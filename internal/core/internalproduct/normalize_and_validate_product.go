package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/pkg/errmsg"
)

func (c *internalProductCore) normalizeAndValidateProduct(
	ctx context.Context,
	data *coreentity.InternalProduct,
) error {
	_ = ctx

	data.Code = strings.ToUpper(strings.TrimSpace(data.Code))
	data.Name = strings.TrimSpace(data.Name)
	data.Description = strings.TrimSpace(data.Description)
	data.Status = strings.ToLower(strings.TrimSpace(data.Status))
	if data.Metadata == nil {
		data.Metadata = map[string]any{}
	}

	if !internalCodePattern.MatchString(data.Code) {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalProductCodeFormatIsInvalid)
	}
	if data.Status != coreentity.InternalProductStatusDraft &&
		data.Status != coreentity.InternalProductStatusActive &&
		data.Status != coreentity.InternalProductStatusInactive &&
		data.Status != coreentity.InternalProductStatusArchived {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalProductStatusIsInvalid)
	}

	return nil
}
