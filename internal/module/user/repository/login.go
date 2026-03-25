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

func (r *userRepo) FindActiveUserByEmailAndTenant(ctx context.Context, email, tenantID string) (*coreentity.User, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.FindActiveUserByEmailAndTenant")
	defer span.End()
	payload := map[string]string{"email": email, "tenantID": tenantID}

	type roleResult struct {
		ID       string `db:"id"`
		RoleName string `db:"role_name"`
	}

	rolesQuery := `
		SELECT ur.role_id as id, r.name as role_name
		FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = (
			SELECT id FROM users WHERE email = ? AND tenant_id = ? AND deleted_at IS NULL
		)
	`

	var roles []roleResult
	err := r.db.SelectContext(ctx, &roles, rolesQuery, email, tenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get roles")
		return nil, err
	}

	if len(roles) == 0 {
		return nil, errmsg.NewCustomErrors(400).SetMessage("User not found or has no roles")
	}

	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.RoleName)
	}

	query := `
		SELECT
			u.id,
			u.email,
			u.password,
			u.tenant_id,
			t.name as tenant_name,
			u.username,
			u.company_id,
			c.name as company_name
		FROM users u
		JOIN tenants t ON t.id = u.tenant_id
		LEFT JOIN org_units c ON c.id = u.company_id AND c.category = 'company'
		WHERE
			u.email = ?
			AND u.tenant_id = ?
			AND u.deleted_at IS NULL
	`

	var row struct {
		ID          string  `db:"id"`
		Email       string  `db:"email"`
		Password    string  `db:"password"`
		TenantID    string  `db:"tenant_id"`
		TenantName  string  `db:"tenant_name"`
		UserName    string  `db:"username"`
		CompanyID   *string `db:"company_id"`
		CompanyName *string `db:"company_name"`
	}

	err = r.db.GetContext(ctx, &row, r.db.Rebind(query), email, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(400).SetMessage("User not found or has no roles")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get user")
		return nil, err
	}

	return &coreentity.User{
		ID:          row.ID,
		Email:       row.Email,
		Password:    row.Password,
		RoleNames:   roleNames,
		TenantID:    row.TenantID,
		TenantName:  row.TenantName,
		UserName:    row.UserName,
		CompanyID:   row.CompanyID,
		CompanyName: row.CompanyName,
	}, nil
}
