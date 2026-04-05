package repository

import (
	"context"

	"codebase-app/internal/entity/common"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r *userRepo) replaceUserRolesTx(ctx context.Context, tx *sqlx.Tx, userID string, roleIDs []string) error {
	deleteQuery := `DELETE FROM user_roles WHERE user_id = ?`
	if _, err := tx.ExecContext(ctx, tx.Rebind(deleteQuery), userID); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"user_id": userID,
		}).Msg("Failed to delete user roles")
		return err
	}

	insertQuery := `INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`
	for _, roleID := range roleIDs {
		if _, err := tx.ExecContext(ctx, tx.Rebind(insertQuery), userID, roleID); err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
				"user_id": userID,
				"role_id": roleID,
			}).Msg("Failed to insert user role")
			return err
		}
	}

	return nil
}
