package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/infrastructure/tracing"
)

func (r *employeeRepo) ExistsByEmployeeNo(ctx context.Context, tenantID, employeeNo, excludeID string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:employee:exists_by_employee_no:ExistsByEmployeeNo")
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
	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &id, exec.Rebind(query), args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
