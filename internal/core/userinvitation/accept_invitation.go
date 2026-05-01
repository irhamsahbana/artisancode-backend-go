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

func (c *userInvitationCore) AcceptInvitation(
	ctx context.Context,
	data coreentity.UserInvitationAcceptPayload,
) (*coreentity.User, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:userinvitation:accept_invitation:AcceptInvitation")
	defer span.End()

	var result *coreentity.User
	err := c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		item, err := c.repo.GetInvitationByTokenHash(txCtx, hashInvitationToken(data.Token))
		if err != nil {
			return err
		}

		item, err = validateAcceptableInvitation(item)
		if err != nil {
			return err
		}

		existingUser, err := c.userRepo.FindActiveUserByEmailAndTenantID(txCtx, item.Email, item.TenantID)
		if err != nil {
			return err
		}
		if existingUser != nil {
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageEmailIsAlreadyRegistered)
		}

		role, err := c.userRepo.GetRoleByName(txCtx, item.RoleCode, item.TenantID)
		if err != nil {
			return err
		}

		userData, err := c.buildAcceptedUser(txCtx, item, data, role.ID)
		if err != nil {
			return err
		}

		userID, err := c.userRepo.InsertUser(txCtx, userData)
		if err != nil {
			log.Ctx(txCtx).Error().Err(err).Any(common.LogKeyPayload, map[string]any{
				"invitation_id": item.ID,
				"email":         item.Email,
				"role_code":     item.RoleCode,
			}).Msg("Failed to insert invited user")
			return err
		}

		if err := c.userRepo.MarkUserEmailVerified(txCtx, userID); err != nil {
			return err
		}

		if item.EmployeeID != nil {
			if err := c.employeeRepo.AssignUser(txCtx, item.TenantID, *item.EmployeeID, userID); err != nil {
				return err
			}
		}

		if err := c.repo.MarkInvitationAccepted(txCtx, item.ID); err != nil {
			return err
		}

		result, err = c.userRepo.GetUser(txCtx, coreentity.User{
			ID:       userID,
			TenantID: item.TenantID,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *userInvitationCore) buildAcceptedUser(
	ctx context.Context,
	item *coreentity.UserInvitation,
	data coreentity.UserInvitationAcceptPayload,
	roleID string,
) (coreentity.User, error) {
	passwordHash, err := hashInvitationPassword(data.Password)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to hash invitation password")
		return coreentity.User{}, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageFailedToAcceptInvitation)
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
			return coreentity.User{}, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageEmployeeAlreadyHasAnActiveUser)
		}

		user.Name = employee.FullName
		user.UserName = employee.EmployeeNo
		return user, nil
	}

	fullName := strings.TrimSpace(data.FullName)
	if fullName == "" {
		return coreentity.User{}, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageFullNameIsRequired)
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
