package repository

import (
	"context"
	"database/sql"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) FindActiveUserByEmailAndTenant(ctx context.Context, email, tenantCode string) (*coreentity.User, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.FindActiveUserByEmailAndTenant")
	defer span.End()

	// Get user base data
	user, err := r.GetUserByEmailAndTenant(ctx, email, tenantCode)
	if err != nil {
		return nil, err
	}

	// Get user roles
	roleNames, err := r.GetUserRolesByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// Attach roles to user object
	user.RoleNames = roleNames

	return user, nil
}

func (r *userRepo) GetUserByEmailAndTenant(ctx context.Context, email, tenantCode string) (*coreentity.User, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetUserByEmailAndTenant")
	defer span.End()

	tenantCode = strings.ToUpper(tenantCode)
	payload := map[string]string{"email": email, "tenantCode": tenantCode}

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
		FROM
			users u
		JOIN tenants t ON t.id = u.tenant_id
		LEFT JOIN org_units c ON c.id = u.company_id AND c.category = 'company'
		WHERE
			u.email = ?
			AND t.code = ?
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

	err := r.db.GetContext(ctx, &row, r.db.Rebind(query), email, tenantCode)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("User not found")
			return nil, errmsg.NewCustomErrors(400).SetMessage("User not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get user details")
		return nil, err
	}

	return &coreentity.User{
		ID:          row.ID,
		Email:       row.Email,
		Password:    row.Password,
		TenantID:    row.TenantID,
		TenantName:  row.TenantName,
		UserName:    row.UserName,
		CompanyID:   row.CompanyID,
		CompanyName: row.CompanyName,
	}, nil
}

func (r *userRepo) GetUserRolesByUserID(ctx context.Context, userID string) ([]string, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetUserRolesByUserID")
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
	err := r.db.SelectContext(ctx, &roles, r.db.Rebind(query), userID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get user roles")
		return nil, err
	}

	if len(roles) == 0 {
		return nil, errmsg.NewCustomErrors(400).SetMessage("User has no roles")
	}

	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.RoleName)
	}

	return roleNames, nil
}
