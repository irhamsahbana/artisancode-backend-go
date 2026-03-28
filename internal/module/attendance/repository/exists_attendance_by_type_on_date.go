package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *attendanceRepo) ExistsAttendanceByTypeOnDate(ctx context.Context, tenantID, employeeID, attendanceDate, attendanceType string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.ExistsAttendanceByTypeOnDate")
	defer span.End()

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM attendance_logs
			WHERE tenant_id = ?
				AND employee_id = ?
				AND attendance_date = ?
				AND type = ?
				AND deleted_at IS NULL
		)
	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, r.db.Rebind(query), tenantID, employeeID, attendanceDate, attendanceType)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"tenant_id":       tenantID,
			"employee_id":     employeeID,
			"attendance_date": attendanceDate,
			"type":            attendanceType,
		}).Msg("Failed to check attendance log existence")
		return false, err
	}

	return exists, nil
}
