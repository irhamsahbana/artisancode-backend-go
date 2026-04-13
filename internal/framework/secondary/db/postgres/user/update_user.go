package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) UpdateUser(ctx context.Context, data coreentity.User) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:update_user:UpdateUser")
	defer span.End()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	args := []any{
		data.Name,
		data.UserName,
		data.Email,
		data.CompanyID,
		data.ID,
		data.TenantID,
	}

	query := `
		UPDATE users
		SET name = ?, username = ?, email = ?, company_id = ?, updated_at = NOW()
	`

	if data.Password != "" {
		query += `, password = ?`
		args = []any{
			data.Name,
			data.UserName,
			data.Email,
			data.CompanyID,
			data.Password,
			data.ID,
			data.TenantID,
		}
	}

	query += `
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	result, err := tx.ExecContext(ctx, tx.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update user")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errmsg.NewCustomErrors(404).SetMessage("User not found")
	}

	if err := r.replaceUserRolesTx(ctx, tx, data.ID, data.RoleIDs); err != nil {
		return err
	}

	return tx.Commit()
}
