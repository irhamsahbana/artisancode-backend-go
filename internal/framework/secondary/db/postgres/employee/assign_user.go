package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *employeeRepo) AssignUser(ctx context.Context, tenantID, employeeID, userID string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:employee:assign_user:AssignUser")
	defer span.End()

	query := `
		UPDATE employees
		SET user_id = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	exec := r.executor(ctx)
	payload := map[string]string{
		"tenant_id":   tenantID,
		"employee_id": employeeID,
		"user_id":     userID,
	}
	result, err := exec.ExecContext(ctx, exec.Rebind(query), userID, employeeID, tenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to assign user to employee")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get assign user rows affected")
		return err
	}
	if rowsAffected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("Employee not found when assigning user")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageEmployeeNotFound)
	}

	return nil
}
