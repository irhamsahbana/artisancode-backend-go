package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *internalUserRepo) CreateInternalUser(
	ctx context.Context,
	data coreentity.InternalUser,
) (*coreentity.InternalUser, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internaluser:repo:CreateInternalUser")
	defer span.End()

	query := `
		INSERT INTO internal_users (
			email,
			full_name,
			password_hash,
			role_code,
			status
		)
		VALUES (
			?,
			?,
			?,
			?,
			?
		)
		RETURNING
			id
	`
	if err := r.db.GetContext(ctx, &data.ID, r.db.Rebind(query),
		data.Email,
		data.FullName,
		data.Password,
		data.RoleCode,
		data.Status,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create internal user")
		return nil, err
	}
	return r.GetInternalUser(ctx, coreentity.InternalUserFilter{ID: data.ID})
}
