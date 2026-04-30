package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func (c *userCore) RegisterOwner(
	ctx context.Context,
	user coreentity.User,
	tenant coreentity.Tenant,
) (*coreentity.RegisterResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:register:RegisterOwner")
	defer span.End()

	payload := map[string]any{
		"user":   user,
		"tenant": tenant,
	}

	var result *coreentity.RegisterResult
	err := c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		tenantExist, err := c.repo.ExistsTenantByCode(txCtx, tenant.Code)
		if err != nil {
			return err
		}
		if tenantExist {
			log.Ctx(txCtx).Warn().Any(common.LogKeyPayload, payload).Msg("Tenant code already registered")
			return errmsg.NewCustomErrors(400).SetMessage("Tenant code is already registered")
		}

		tenantID, _, err := c.createTenant(txCtx, tenant)
		if err != nil {
			return err
		}

		ownerRole, err := c.getOwnerRole(txCtx, tenantID)
		if err != nil {
			return err
		}

		userID, err := c.createOwnerUser(txCtx, user, tenantID, ownerRole.ID)
		if err != nil {
			return err
		}

		err = c.issueEmailVerification(txCtx, coreentity.User{
			ID:                userID,
			Name:              user.Name,
			Email:             user.Email,
			TenantID:          tenantID,
			TenantName:        user.TenantName,
			PreferredLanguage: tenant.PreferredLanguage,
		})
		if err != nil {
			return err
		}

		result = &coreentity.RegisterResult{
			Email:                user.Email,
			VerificationRequired: true,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *userCore) createTenant(ctx context.Context, tenant coreentity.Tenant) (string, string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:register:createTenant")
	defer span.End()

	preferredLanguage := tenant.PreferredLanguage
	if preferredLanguage == "" {
		preferredLanguage = "id"
	}

	tenantID, err := c.repo.InsertTenant(ctx, coreentity.Tenant{
		Name: tenant.Name,
		Code: tenant.Code,
	})
	if err != nil {
		return "", "", err
	}

	companyID, err := c.repo.InitializeTenant(ctx, tenantID, tenant.Name, preferredLanguage)
	if err != nil {
		return "", "", err
	}

	return tenantID, companyID, nil
}

func (c *userCore) getOwnerRole(ctx context.Context, tenantID string) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:register:getOwnerRole")
	defer span.End()

	return c.repo.GetRoleByName(ctx, "owner", tenantID)
}

func (c *userCore) createOwnerUser(ctx context.Context, user coreentity.User, tenantID, roleID string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:register:createOwnerUser")
	defer span.End()

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	userData := coreentity.User{
		RoleIDs:    []string{roleID},
		RoleNames:  []string{"owner"},
		Name:       user.Name,
		UserName:   user.UserName,
		Email:      user.Email,
		Password:   string(hashed),
		TenantID:   tenantID,
		TenantName: user.TenantName,
	}

	return c.repo.InsertUser(ctx, userData)
}
