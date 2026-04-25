package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *userRepo) InsertUser(ctx context.Context, user coreentity.User) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:register_owner:InsertUser")
	defer span.End()

	query := `
		INSERT INTO users (
			tenant_id, company_id, name, username, email, password
		) VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`

	var userID string
	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &userID, exec.Rebind(query),
		user.TenantID,
		user.CompanyID,
		user.Name,
		user.UserName,
		user.Email,
		user.Password,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, user).Msg("Failed to insert user")
		return "", err
	}

	roleQuery := `INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`
	for _, roleID := range user.RoleIDs {
		_, err = exec.ExecContext(ctx, exec.Rebind(roleQuery), userID, roleID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"userID": userID, "roleID": roleID}).Msg("Failed to insert user_role")
			return "", err
		}
	}

	return userID, nil
}
