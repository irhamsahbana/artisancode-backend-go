package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) GetInvoices(
	ctx context.Context,
	filter coreentity.TenantBillingInvoiceListFilter,
) ([]coreentity.TenantBillingInvoice, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:get_invoices:GetInvoices",
	)
	defer span.End()

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Paginate < 1 {
		filter.Paginate = 20
	}
	if filter.Paginate > 100 {
		filter.Paginate = 100
	}

	items := make([]coreentity.TenantBillingInvoice, 0)
	query := `
		WITH latest_attempt AS (
			SELECT DISTINCT ON (internal_tenant_invoice_id)
				id,
				internal_tenant_invoice_id,
				provider_payment_url
			FROM internal_tenant_payment_attempts
			WHERE tenant_id = ?
				AND deleted_at IS NULL
				AND status IN ('initiated', 'pending')
			ORDER BY internal_tenant_invoice_id, created_at DESC
		)
		SELECT
			i.id,
			i.invoice_number,
			i.status,
			i.amount::text AS amount,
			i.currency_code AS currency,
			a.id AS payment_attempt_id,
			a.provider_payment_url AS payment_url,
			i.due_at::text AS due_at,
			i.paid_at::text AS paid_at,
			i.created_at::text AS created_at
		FROM internal_tenant_invoices i
		LEFT JOIN latest_attempt a ON a.internal_tenant_invoice_id = i.id
		WHERE i.tenant_id = ?
			AND i.deleted_at IS NULL
		ORDER BY i.created_at DESC
		LIMIT ?
		OFFSET ?
	`
	if err := r.exec(ctx).SelectContext(
		ctx,
		&items,
		r.exec(ctx).Rebind(query),
		filter.UserCtx.TenantID,
		filter.UserCtx.TenantID,
		filter.Paginate,
		(filter.Page-1)*filter.Paginate,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get tenant billing invoices")
		return nil, err
	}

	return items, nil
}
