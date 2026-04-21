package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func (c *userInvitationCore) AcceptInvitation(ctx context.Context, data coreentity.UserInvitationAcceptPayload) (*coreentity.User, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:userinvitation:accept_invitation:AcceptInvitation")
	defer span.End()

	item, err := c.repo.GetInvitationByTokenHash(ctx, hashInvitationToken(data.Token))
	if err != nil {
		return nil, err
	}

	item, err = validateAcceptableInvitation(item)
	if err != nil {
		return nil, err
	}

	existingUser, err := c.userRepo.FindActiveUserByEmailAndTenantID(ctx, item.Email, item.TenantID)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Email is already registered")
	}

	role, err := c.userRepo.GetRoleByName(ctx, item.RoleCode, item.TenantID)
	if err != nil {
		return nil, err
	}

	userData, err := c.buildAcceptedUser(ctx, item, data, role.ID)
	if err != nil {
		return nil, err
	}

	userID, err := c.userRepo.InsertUser(ctx, userData)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]any{
			"invitation_id": item.ID,
			"email":         item.Email,
			"role_code":     item.RoleCode,
		}).Msg("Failed to insert invited user")
		return nil, err
	}

	if err := c.userRepo.MarkUserEmailVerified(ctx, userID); err != nil {
		return nil, err
	}

	if item.EmployeeID != nil {
		if err := c.employeeRepo.AssignUser(ctx, item.TenantID, *item.EmployeeID, userID); err != nil {
			return nil, err
		}
	}

	if err := c.repo.MarkInvitationAccepted(ctx, item.ID); err != nil {
		return nil, err
	}

	return c.userRepo.GetUser(ctx, coreentity.User{
		ID:       userID,
		TenantID: item.TenantID,
	})
}

func (c *userInvitationCore) buildAcceptedUser(ctx context.Context, item *coreentity.UserInvitation, data coreentity.UserInvitationAcceptPayload, roleID string) (coreentity.User, error) {
	passwordHash, err := hashInvitationPassword(data.Password)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to hash invitation password")
		return coreentity.User{}, errmsg.NewCustomErrors(500).SetMessage("Failed to accept invitation")
	}

	user := coreentity.User{
		RoleIDs:   []string{roleID},
		RoleNames: []string{item.RoleCode},
		Email:     item.Email,
		Password:  passwordHash,
		TenantID:  item.TenantID,
	}

	if item.EmployeeID != nil {
		employee, err := c.employeeRepo.GetEmployee(ctx, coreentity.Employee{
			TenantID: item.TenantID,
			ID:       *item.EmployeeID,
		})
		if err != nil {
			return coreentity.User{}, err
		}
		if employee.UserID != nil {
			return coreentity.User{}, errmsg.NewCustomErrors(400).SetMessage("Employee already has an active user")
		}

		user.Name = employee.FullName
		user.UserName = employee.EmployeeNo
		return user, nil
	}

	fullName := strings.TrimSpace(data.FullName)
	if fullName == "" {
		return coreentity.User{}, errmsg.NewCustomErrors(400).SetMessage("Full name is required")
	}

	user.Name = fullName
	user.UserName = deriveUsernameFromEmail(item.Email)

	return user, nil
}

func hashInvitationPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
