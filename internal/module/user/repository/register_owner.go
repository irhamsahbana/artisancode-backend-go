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

func (r *userRepo) ExistsActiveUserByEmail(ctx context.Context, email string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.ExistsActiveUserByEmail")
	defer span.End()

	var existing string
	query := `SELECT id FROM users WHERE email = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &existing, r.db.Rebind(query), email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, email).Msg("User not found")
			return false, nil
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, email).Msg("Failed to check user existence")
		return false, err
	}
	return existing != "", nil
}

func (r *userRepo) GetRoleByName(ctx context.Context, roleName, tenantID string) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetRoleByName")
	defer span.End()

	type dao struct {
		ID       string  `db:"id"`
		TenantID *string `db:"tenant_id"`
		Name     string  `db:"name"`
	}
	var row dao

	query := `
		SELECT id, tenant_id, name
		FROM roles
		WHERE
			name = ?
			AND (tenant_id = ? OR tenant_id IS NULL)
			AND deleted_at IS NULL
		ORDER BY CASE WHEN tenant_id = ? THEN 0 ELSE 1 END
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &row, r.db.Rebind(query), roleName, tenantID, tenantID)
	if err != nil {
		payload := map[string]string{"roleName": roleName, "tenantID": tenantID}
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("Role not found")
			return nil, errmsg.NewCustomErrors(400).SetMessage("Role not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get role")
		return nil, err
	}

	return &coreentity.Role{
		ID:       row.ID,
		TenantID: row.TenantID,
		Name:     row.Name,
	}, nil
}

func (r *userRepo) InsertUser(ctx context.Context, user coreentity.User) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.InsertUser")
	defer span.End()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO users (
			tenant_id, company_id, name, username, email, password
		) VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`

	var userID string
	err = tx.GetContext(ctx, &userID, tx.Rebind(query),
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
		_, err = tx.ExecContext(ctx, tx.Rebind(roleQuery), userID, roleID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"userID": userID, "roleID": roleID}).Msg("Failed to insert user_role")
			return "", err
		}
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return userID, nil
}
