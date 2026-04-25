package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	postgresTx "codebase-app/internal/framework/secondary/db/postgres/transaction"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) replaceUserRoles(ctx context.Context, exec postgresTx.SQLExecutor, userID string, roleIDs []string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:replace_user_roles_tx:replaceUserRolesTx")
	defer span.End()

	deleteQuery := `DELETE FROM user_roles WHERE user_id = ?`
	if _, err := exec.ExecContext(ctx, exec.Rebind(deleteQuery), userID); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"user_id": userID,
		}).Msg("Failed to delete user roles")
		return err
	}

	insertQuery := `INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`
	for _, roleID := range roleIDs {
		if _, err := exec.ExecContext(ctx, exec.Rebind(insertQuery), userID, roleID); err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
				"user_id": userID,
				"role_id": roleID,
			}).Msg("Failed to insert user role")
			return err
		}
	}

	return nil
}
