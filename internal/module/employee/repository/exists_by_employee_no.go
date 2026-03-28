package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/infrastructure/tracing"
)

func (r *employeeRepo) ExistsByEmployeeNo(ctx context.Context, tenantID, employeeNo, excludeID string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.ExistsByEmployeeNo")
	defer span.End()

	query := `
		SELECT id FROM employees
		WHERE tenant_id = ? AND employee_no = ? AND deleted_at IS NULL
	`
	args := []any{tenantID, employeeNo}

	if excludeID != "" {
		query += ` AND id != ?`
		args = append(args, excludeID)
	}

	var id string
	err := r.db.GetContext(ctx, &id, r.db.Rebind(query), args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}