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

	err := normalizeEmployeeJoinDate(&data)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, map[string]any{
			"join_date":          data.JoinDate,
			"join_date_timezone": data.JoinDateTimezone,
		}).Msg("Invalid join date payload")
		return errmsg.NewCustomErrors(400).SetMessage("Invalid join date or join date timezone")
	}

	existing, err := c.repo.GetEmployee(ctx, coreentity.Employee{
		TenantID: data.TenantID,
		ID:       data.ID,
	})
	if err != nil {
		return err
	}

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

	emailChanged := existing.Email != data.Email
	if emailChanged {
		emailExists, err := c.userRepo.ExistsActiveUserByEmail(ctx, data.Email)
		if err != nil {
			return err
		}
		if emailExists {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
				"email": data.Email,
			}).Msg("Email already registered")
			return errmsg.NewCustomErrors(400).SetMessage("Email is already registered")
		}
	}

	data.UserID = existing.UserID

	err = c.repo.UpdateEmployee(ctx, data)
	if err != nil {
		return err
	}

	if emailChanged && existing.UserID != nil {
		err = c.userRepo.UpdateUserEmail(ctx, *existing.UserID, data.TenantID, data.Email)
		if err != nil {
			return err
		}
	}

	if data.Password != "" && existing.UserID != nil {
		hashedPassword, err := hashPassword(data.Password)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to hash password")
			return errmsg.NewCustomErrors(500).SetMessage("Failed to update employee password")
		}

		err = c.userRepo.UpdateUserPassword(ctx, *existing.UserID, data.TenantID, hashedPassword)
		if err != nil {
			return err
		}
	}

	return nil
}
