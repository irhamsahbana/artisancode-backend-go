package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) insertDefaultHeadquarterBranch(
	ctx context.Context,
	tx sqlExecutor,
	tenantID string,
	companyID string,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:initialize_tenant:insertDefaultHeadquarterBranch",
	)
	defer span.End()

	query := `
		INSERT INTO org_units (tenant_id, parent_id, name, code, category, config)
		VALUES (?, ?, 'Headquarter', 'HEADQUARTER', 'branch', '{}')
		ON CONFLICT (tenant_id, code) DO UPDATE SET
			parent_id = EXCLUDED.parent_id,
			name = EXCLUDED.name,
			category = EXCLUDED.category,
			updated_at = CURRENT_TIMESTAMP,
			deleted_at = NULL
	`
	if _, err := tx.ExecContext(ctx, tx.Rebind(query), tenantID, companyID); err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID, "companyID": companyID}).
			Msg("Failed to insert default headquarter branch")
		return err
	}

	return nil
}
