package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

type googleTenantSignupInput struct {
	TenantName        string
	TenantCode        string
	PreferredLanguage string
}

func (c *userCore) registerGoogleTenantOwner(
	ctx context.Context,
	identity *coreentity.GoogleIdentity,
	input googleTenantSignupInput,
) (*coreentity.GoogleRegisterResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:tenant_signup:registerGoogleTenantOwner")
	defer span.End()

	var result *coreentity.GoogleRegisterResult
	err := c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := c.ensureGoogleSignupTenantCodeAvailable(txCtx, input.TenantCode); err != nil {
			return err
		}

		if err := c.ensureGoogleIdentityNotLinked(txCtx, identity); err != nil {
			return err
		}

		if err := c.ensureGoogleSignupEmailAvailable(txCtx, identity.Email); err != nil {
			return err
		}

		tenantID, companyID, err := c.createTenant(txCtx, coreentity.Tenant{
			Name:              input.TenantName,
			Code:              input.TenantCode,
			PreferredLanguage: input.PreferredLanguage,
		})
		if err != nil {
			return err
		}

		ownerRole, err := c.getOwnerRole(txCtx, tenantID)
		if err != nil {
			return err
		}

		user, err := c.createGoogleTenantOwnerUser(txCtx, identity, googleTenantOwnerInput{
			TenantID:          tenantID,
			TenantName:        input.TenantName,
			CompanyID:         companyID,
			CompanyName:       input.TenantName,
			RoleID:            ownerRole.ID,
			PreferredLanguage: input.PreferredLanguage,
		})
		if err != nil {
			return err
		}

		if err := c.createGoogleAuthIdentity(txCtx, tenantID, user.ID, identity); err != nil {
			return err
		}

		tokens, err := c.issueAuthTokens(txCtx, user)
		if err != nil {
			return err
		}

		result = &coreentity.GoogleRegisterResult{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			TenantCode:   input.TenantCode,
		}
		return nil
	})
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	return result, nil
}

func (c *userCore) ensureGoogleSignupTenantCodeAvailable(
	ctx context.Context,
	tenantCode string,
) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:tenant_signup:ensureGoogleSignupTenantCodeAvailable")
	defer span.End()

	tenantExist, err := c.repo.ExistsTenantByCode(ctx, tenantCode)
	if err != nil {
		return err
	}
	if tenantExist {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"tenant_code": tenantCode,
		}).Msg("Tenant code already used")
		return codedError(400, errmsg.MessageTenantCodeIsAlreadyRegistered, errorCodeTenantCodeAlreadyUsed)
	}

	return nil
}

func (c *userCore) ensureGoogleIdentityNotLinked(
	ctx context.Context,
	identity *coreentity.GoogleIdentity,
) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:tenant_signup:ensureGoogleIdentityNotLinked")
	defer span.End()

	linkedIdentity, err := c.repo.FindAuthIdentityByProviderSubject(
		ctx,
		coreentity.AuthProviderGoogle,
		identity.Subject,
	)
	if err != nil {
		return err
	}
	if linkedIdentity != nil {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"provider": coreentity.AuthProviderGoogle,
		}).Msg("Google identity already linked")
		return codedError(
			400,
			errmsg.MessageGoogleIdentityIsAlreadyLinked,
			errorCodeGoogleIdentityAlreadyLinked,
		)
	}

	return nil
}

func (c *userCore) ensureGoogleSignupEmailAvailable(
	ctx context.Context,
	email string,
) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:tenant_signup:ensureGoogleSignupEmailAvailable")
	defer span.End()

	users, err := c.repo.FindActiveUsersByEmail(ctx, email)
	if err != nil {
		return err
	}
	if len(users) > 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"email": email,
		}).Msg("Google register email already registered")
		return codedError(
			400,
			errmsg.MessageGoogleEmailIsAlreadyRegistered,
			errorCodeGoogleEmailAlreadyRegistered,
		)
	}

	return nil
}

type googleTenantOwnerInput struct {
	TenantID          string
	TenantName        string
	CompanyID         string
	CompanyName       string
	RoleID            string
	PreferredLanguage string
}

func (c *userCore) createGoogleTenantOwnerUser(
	ctx context.Context,
	identity *coreentity.GoogleIdentity,
	input googleTenantOwnerInput,
) (coreentity.User, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:tenant_signup:createGoogleTenantOwnerUser")
	defer span.End()

	now := time.Now().UTC()
	user := coreentity.User{
		RoleIDs:           []string{input.RoleID},
		RoleNames:         []string{"owner"},
		Name:              displayName(identity),
		UserName:          usernameFromEmail(identity.Email),
		Email:             identity.Email,
		TenantID:          input.TenantID,
		TenantName:        input.TenantName,
		CompanyID:         &input.CompanyID,
		CompanyName:       &input.CompanyName,
		EmailVerifiedAt:   &now,
		PreferredLanguage: input.PreferredLanguage,
	}

	userID, err := c.repo.InsertUser(ctx, user)
	if err != nil {
		return coreentity.User{}, err
	}
	user.ID = userID

	return user, nil
}

func (c *userCore) createGoogleAuthIdentity(
	ctx context.Context,
	tenantID string,
	userID string,
	identity *coreentity.GoogleIdentity,
) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:tenant_signup:createGoogleAuthIdentity")
	defer span.End()

	return c.repo.CreateAuthIdentity(ctx, coreentity.UserAuthIdentity{
		TenantID:        tenantID,
		UserID:          userID,
		Provider:        coreentity.AuthProviderGoogle,
		ProviderSubject: identity.Subject,
		Email:           identity.Email,
		EmailVerified:   identity.EmailVerified,
		DisplayName:     identity.DisplayName,
		PictureURL:      identity.PictureURL,
	})
}
