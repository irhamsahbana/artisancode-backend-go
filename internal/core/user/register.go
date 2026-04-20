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

func (c *userCore) RegisterOwner(ctx context.Context, user coreentity.User, tenant coreentity.Tenant) (*coreentity.RegisterResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:register:RegisterOwner")
	defer span.End()

	exist, err := c.repo.ExistsActiveUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if exist {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{"email": user.Email}).Msg("Email already registered")
		return nil, errmsg.NewCustomErrors(400).SetMessage("Email is already registered")
	}

	tenantExist, err := c.repo.ExistsTenantByCode(ctx, tenant.Code)
	if err != nil {
		return nil, err
	}
	if tenantExist {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{"tenantCode": tenant.Code}).Msg("Tenant code already registered")
		return nil, errmsg.NewCustomErrors(400).SetMessage("Tenant code is already registered")
	}

	tenantID, err := c.createTenant(ctx, tenant)
	if err != nil {
		return nil, err
	}

	ownerRole, err := c.getOwnerRole(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	userID, err := c.createOwnerUser(ctx, user, tenantID, ownerRole.ID)
	if err != nil {
		return nil, err
	}

	err = c.issueEmailVerification(ctx, coreentity.User{
		ID:                userID,
		Name:              user.Name,
		Email:             user.Email,
		TenantID:          tenantID,
		TenantName:        user.TenantName,
		PreferredLanguage: tenant.PreferredLanguage,
	})
	if err != nil {
		return nil, err
	}

	return &coreentity.RegisterResult{
		Email:                user.Email,
		VerificationRequired: true,
	}, nil
}

func (c *userCore) createTenant(ctx context.Context, tenant coreentity.Tenant) (string, error) {
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
		return "", err
	}

	_, err = c.repo.InitializeTenant(ctx, tenantID, tenant.Name, preferredLanguage)
	if err != nil {
		return "", err
	}

	return tenantID, nil
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
