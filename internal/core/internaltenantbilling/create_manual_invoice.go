package core

import (
	"context"
	"fmt"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalTenantBillingCore) CreateManualInvoice(
	ctx context.Context,
	input coreentity.InternalBillingManualInvoiceInput,
) (*coreentity.InternalBillingManualInvoiceResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:create_manual_invoice:CreateManualInvoice")
	defer span.End()

	if input.TenantID == "" {
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("tenant_id is required"))
	}
	if input.Amount == "" {
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("amount is required"))
	}
	if input.Currency == "" {
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("currency is required"))
	}

	now := time.Now()
	invNumber := fmt.Sprintf("INV-MAN-%s-%d", input.TenantID[:min(8, len(input.TenantID))], now.Unix())

	var items []map[string]interface{}
	for _, it := range input.Items {
		items = append(items, map[string]interface{}{
			"description": it.Description,
			"amount":      it.Amount,
			"quantity":    it.Quantity,
		})
	}

	dueAt := (*string)(nil)
	if input.DueAt != "" {
		dueAt = &input.DueAt
	}

	invoice := coreentity.InternalTenantInvoice{
		TenantID:                     input.TenantID,
		InternalTenantSubscriptionID:  nil,
		InvoiceNumber:                invNumber,
		Status:                       coreentity.InternalTenantInvoiceStatusPending,
		CurrencyCode:                 input.Currency,
		Amount:                       input.Amount,
		AmountPaid:                   "0",
		AmountOutstanding:            input.Amount,
		SourceType:                   coreentity.InternalBillingSourceTypeManualInvoice,
		TargetSubscriptionState:      "",
		DueAt:                        dueAt,
		Metadata:                     map[string]interface{}{"items": items, "notes": input.Description},
	}

	created, err := c.billingRepo.CreateInternalInvoice(ctx, invoice)
	if err != nil {
		return nil, err
	}

	return &coreentity.InternalBillingManualInvoiceResult{
		Invoice: *created,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
