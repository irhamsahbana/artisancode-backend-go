package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"github.com/rs/zerolog/log"
)

func (r *userRepo) GetUserRolesByUserID(ctx context.Context, userID string) ([]string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:login:GetUserRolesByUserID")
	defer span.End()

	payload := map[string]string{"userID": userID}

	type roleResult struct {
		RoleName string `db:"role_name"`
	}

	query := `
		SELECT r.name as role_name
		FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = ?
	`

	var roles []roleResult
	exec := r.executor(ctx)
	err := exec.SelectContext(ctx, &roles, exec.Rebind(query), userID)
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get user roles")
		return nil, err
	}

	if len(roles) == 0 {
		err = errmsg.NewCustomErrors(400).SetMessage("User has no roles")
		tracing.RecordError(span, err)
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("User has no roles")
		return nil, err
	}

	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.RoleName)
	}

	return roleNames, nil
}
