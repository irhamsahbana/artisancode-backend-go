package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (r *internalQuotationRepo) MarkQuotationConverted(
	ctx context.Context,
	tenantID string,
	quotationID string,
	orderID string,
) (*coreentity.InternalQuotation, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalquotation:mark_quotation_converted:MarkQuotationConverted",
	)
	defer span.End()

	query := `
		UPDATE internal_quotations
		SET
			status = 'converted',
			converted_to_order_id = ?,
			updated_at = NOW()
		WHERE id = ?
			AND tenant_id = ?
			AND deleted_at IS NULL
	`
	result, err := r.exec(ctx).ExecContext(ctx, r.exec(ctx).Rebind(query), orderID, quotationID, tenantID)
	if err != nil {
		return nil, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return nil, errmsg.NewCustomErrors(404).SetMessage("Quotation not found")
	}
	return r.GetQuotation(ctx, tenantID, quotationID)
}
