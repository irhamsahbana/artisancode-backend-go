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

func (r *internalUserRepo) GetInternalUser(ctx context.Context, filter coreentity.InternalUserFilter) (*coreentity.InternalUser, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internaluser:repo:GetInternalUser")
	defer span.End()

	query := `
		SELECT id, full_name, email, password_hash, role_code, status,
		       last_login_at::text AS last_login_at, created_at::text AS created_at, updated_at::text AS updated_at
		FROM internal_users
		WHERE id = ? AND deleted_at IS NULL
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

	if err := r.db.GetContext(ctx, &row, r.db.Rebind(query), filter.ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Internal user not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get internal user")
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
