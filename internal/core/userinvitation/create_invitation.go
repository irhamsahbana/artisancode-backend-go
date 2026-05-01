package core

import (
	"context"
	"strings"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *userInvitationCore) CreateInvitation(
	ctx context.Context,
	data coreentity.UserInvitation,
) (*coreentity.UserInvitation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:userinvitation:create_invitation:CreateInvitation")
	defer span.End()

	if !data.UserCtx.HasRole("owner") && !data.UserCtx.HasRole("admin") {
		return nil, errmsg.NewCustomErrors(403).SetMessage(errmsg.MessageYouAreNotAuthorizedToInviteUsers)
	}

	data.Email = normalizeInvitationEmail(data.Email)
	data.RoleCode = strings.TrimSpace(strings.ToLower(data.RoleCode))
	data.TenantID = data.UserCtx.TenantID
	data.InvitedBy = data.UserCtx.UserID
	data.Status = coreentity.UserInvitationStatusPending
	data.LastSentAt = time.Now()
	data.ExpiresAt = data.LastSentAt.Add(invitationExpiryDuration)

	if err := c.validateInvitationTarget(ctx, data); err != nil {
		return nil, err
	}

	rawToken, tokenHash, err := generateInvitationToken()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to generate invitation token")
		return nil, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageFailedToCreateInvitation)
	}

	data.AcceptToken = rawToken
	data.TokenHash = tokenHash

	var created *coreentity.UserInvitation
	err = c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		exists, err := c.repo.ExistsActiveInvitation(txCtx, data.TenantID, data.Email, data.RoleCode, data.EmployeeID)
		if err != nil {
			return err
		}
		if exists {
			log.Ctx(txCtx).Warn().Any(common.LogKeyPayload, data).Msg("Active invitation already exists")
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageAnActiveInvitationAlreadyExists)
		}

		created, err = c.repo.CreateInvitation(txCtx, data)
		if err != nil {
			return err
		}

		created.AcceptToken = rawToken
		emailSent, err := c.trySendInvitationEmail(txCtx, created)
		if err != nil {
			return err
		}
		created.EmailSent = emailSent
		return nil
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (c *userInvitationCore) validateInvitationTarget(ctx context.Context, data coreentity.UserInvitation) error {
	if data.RoleCode != coreentity.UserInvitationRoleAdmin && data.RoleCode != coreentity.UserInvitationRoleEmployee {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInvalidInvitationRole)
	}

	if data.RoleCode == coreentity.UserInvitationRoleAdmin {
		if !data.UserCtx.HasRole("owner") {
			return errmsg.NewCustomErrors(403).SetMessage(errmsg.MessageOnlyOwnerCanInviteAdminUsers)
		}
		if data.EmployeeID != nil {
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageEmployeeIdIsNotAllowedForAdminInvitations)
		}
	}

	if data.RoleCode == coreentity.UserInvitationRoleEmployee {
		if data.EmployeeID == nil || *data.EmployeeID == "" {
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageEmployeeIdIsRequiredForEmployeeInvitations)
		}

		employee, err := c.employeeRepo.GetEmployee(ctx, coreentity.Employee{
			TenantID: data.TenantID,
			ID:       *data.EmployeeID,
		})
		if err != nil {
			return err
		}
		if employee.UserID != nil {
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageEmployeeAlreadyHasAnActiveUser)
		}
		if normalizeInvitationEmail(employee.Email) != data.Email {
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInvitationEmailMustMatchEmployeeEmail)
		}
	}

	existingUser, err := c.userRepo.FindActiveUserByEmailAndTenantID(ctx, data.Email, data.TenantID)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageEmailIsAlreadyRegistered)
	}

	if _, err := c.userRepo.GetRoleByName(ctx, data.RoleCode, data.TenantID); err != nil {
		return err
	}

	return nil
}
