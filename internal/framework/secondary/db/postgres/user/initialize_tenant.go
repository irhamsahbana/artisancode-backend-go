package repository

import (
	"codebase-app/internal/infrastructure/tracing"
	"context"
)

func (r *userRepo) InitializeTenant(
	ctx context.Context,
	tenantID string,
	companyName string,
	preferredLanguage string,
) (string, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:initialize_tenant:InitializeTenant",
	)
	defer span.End()

	tx := r.executor(ctx)

	companyID, err := r.insertDefaultCompany(ctx, tx, tenantID, companyName, preferredLanguage)
	if err != nil {
		return "", err
	}

	err = r.insertDefaultHeadquarterBranch(ctx, tx, tenantID, companyID)
	if err != nil {
		return "", err
	}

	err = r.copyTemplateRoles(ctx, tx, tenantID)
	if err != nil {
		return "", err
	}

	err = r.copyTemplatePermissions(ctx, tx, tenantID)
	if err != nil {
		return "", err
	}

	err = r.copyTemplateRolePermissions(ctx, tx, tenantID)
	if err != nil {
		return "", err
	}

	return companyID, nil
}
