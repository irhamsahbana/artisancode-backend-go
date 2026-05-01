package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *employeeCore) CreateEmployee(ctx context.Context, data coreentity.Employee) (*coreentity.Employee, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:employee:create_employee:CreateEmployee")
	defer span.End()

	err := normalizeEmployeeJoinDate(&data)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, map[string]any{
			"join_date":          data.JoinDate,
			"join_date_timezone": data.JoinDateTimezone,
		}).Msg("Invalid join date payload")
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInvalidJoinDateOrJoinDateTimezone)
	}
	if data.ShiftID == nil || *data.ShiftID == "" {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"employee_no": data.EmployeeNo,
			"tenant_id":   data.TenantID,
		}).Msg("Work shift is required when creating employee")
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageWorkShiftIsRequired)
	}

	// Validate unique employee_no per tenant
	exists, err := c.repo.ExistsByEmployeeNo(ctx, data.TenantID, data.EmployeeNo, "")
	if err != nil {
		return nil, err
	}
	if exists {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"employee_no": data.EmployeeNo,
			"tenant_id":   data.TenantID,
		}).Msg("Employee number already exists")
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageEmployeeNumberAlreadyExistsInThisTenant)
	}

	// Keep employee creation independent from login access.
	// User accounts are created later through the invitation flow.
	data.UserID = nil

	return c.repo.CreateEmployee(ctx, data)
}
