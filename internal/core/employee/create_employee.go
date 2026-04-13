package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg"
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
		return nil, errmsg.NewCustomErrors(400).SetMessage("Invalid join date or join date timezone")
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
		return nil, errmsg.NewCustomErrors(400).SetMessage("Employee number already exists in this tenant")
	}

	// Check if email is already used by another user
	emailExists, err := c.userRepo.ExistsActiveUserByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"email": data.Email,
		}).Msg("Email already registered")
		return nil, errmsg.NewCustomErrors(400).SetMessage("Email is already registered")
	}

	// Create user account for the employee
	userID, err := c.createEmployeeUser(ctx, data)
	if err != nil {
		return nil, err
	}

	data.UserID = &userID

	// Create the employee record
	return c.repo.CreateEmployee(ctx, data)
}

// createEmployeeUser creates a user account with employee role
func (c *employeeCore) createEmployeeUser(ctx context.Context, data coreentity.Employee) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:employee:create_employee:createEmployeeUser")
	defer span.End()

	// Get the employee role
	role, err := c.userRepo.GetRoleByName(ctx, "employee", data.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get employee role")
		return "", errmsg.NewCustomErrors(500).SetMessage("Failed to get employee role")
	}

	// Generate random password
	plainPassword := data.Password
	if plainPassword == "" {
		plainPassword = pkg.GeneratePassword(12)
	}

	hashed, err := hashPassword(plainPassword)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to hash password")
		return "", errmsg.NewCustomErrors(500).SetMessage("Failed to create user account")
	}

	user := coreentity.User{
		RoleIDs:   []string{role.ID},
		RoleNames: []string{"employee"},
		Name:      data.FullName,
		UserName:  data.EmployeeNo,
		Email:     data.Email,
		Password:  hashed,
		TenantID:  data.TenantID,
	}

	userID, err := c.userRepo.InsertUser(ctx, user)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, user).Msg("Failed to insert user for employee")
		return "", err
	}

	return userID, nil
}
