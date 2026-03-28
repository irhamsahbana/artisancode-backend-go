package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *employeeCore) UpdateEmployee(ctx context.Context, data coreentity.Employee) error {
	ctx, span := tracing.StartSpan(ctx, "core.UpdateEmployee")
	defer span.End()

	// Validate unique employee_no per tenant (exclude self)
	exists, err := c.repo.ExistsByEmployeeNo(ctx, data.TenantID, data.EmployeeNo, data.ID)
	if err != nil {
		return err
	}
	if exists {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"employee_no": data.EmployeeNo,
			"tenant_id":   data.TenantID,
		}).Msg("Employee number already exists")
		return errmsg.NewCustomErrors(400).SetMessage("Employee number already exists in this tenant")
	}

	return c.repo.UpdateEmployee(ctx, data)
}