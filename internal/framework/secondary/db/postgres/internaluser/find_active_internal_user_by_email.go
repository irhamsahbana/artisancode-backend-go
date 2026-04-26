package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"github.com/rs/zerolog/log"
)

func (r *internalUserRepo) FindActiveInternalUserByEmail(
	ctx context.Context,
	email string,
) (*coreentity.InternalUser, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaluser:repo:FindActiveInternalUserByEmail",
	)
	defer span.End()

	query := `
		SELECT id, full_name, email, password_hash, role_code, status,
		       last_login_at::text AS last_login_at, created_at::text AS created_at, updated_at::text AS updated_at
		FROM internal_users
		WHERE LOWER(email) = LOWER(?) AND deleted_at IS NULL
	`

	var row struct {
		ID          string  `db:"id"`
		FullName    string  `db:"full_name"`
		Email       string  `db:"email"`
		Password    string  `db:"password_hash"`
		RoleCode    string  `db:"role_code"`
		Status      string  `db:"status"`
		LastLoginAt *string `db:"last_login_at"`
		CreatedAt   string  `db:"created_at"`
		UpdatedAt   *string `db:"updated_at"`
	}

	if err := r.db.GetContext(ctx, &row, r.db.Rebind(query), email); err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, email).Msg("Internal user not found by email")
			return nil, errmsg.NewCustomErrors(400).SetMessage("User not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, email).Msg("Failed to find internal user by email")
		return nil, err
	}

	item := coreentity.InternalUser{
		ID:          row.ID,
		FullName:    row.FullName,
		Email:       row.Email,
		Password:    row.Password,
		RoleCode:    row.RoleCode,
		Status:      row.Status,
		LastLoginAt: row.LastLoginAt,
		CreatedAt:   row.CreatedAt,
	}
	if row.UpdatedAt != nil {
		item.UpdatedAt = *row.UpdatedAt
	}
	return &item, nil
}
