package core

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

var (
	tenantCodePattern  = regexp.MustCompile(`^[A-HJ-NP-Z2-9]{3,5}$`)
	tenantCodeReserved = map[string]struct{}{
		"ADMIN": {},
		"OWNER": {},
		"LOGIN": {},
		"ROOT":  {},
		"TEST":  {},
		"NULL":  {},
		"API":   {},
		"APP":   {},
		"WWW":   {},
	}
)

const (
	errorCodeGoogleAccountNotConnected       = "google_account_not_connected"
	errorCodeGoogleEmailAmbiguous            = "google_email_ambiguous"
	errorCodeGoogleEmailAlreadyRegistered    = "google_email_already_registered"
	errorCodeGoogleIdentityAlreadyLinked     = "google_identity_already_linked"
	errorCodeTenantCodeInvalid               = "tenant_code_invalid"
	errorCodeTenantCodeReserved              = "tenant_code_reserved"
	errorCodeTenantCodeAlreadyUsed           = "tenant_code_already_used"
	errorCodeTenantSetupConfirmationRequired = "tenant_setup_confirmation_required"
)

func (c *userCore) GoogleRegister(
	ctx context.Context,
	input coreentity.GoogleRegisterInput,
) (*coreentity.GoogleRegisterResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:google_auth:GoogleRegister")
	defer span.End()

	if !input.ConfirmTenantSetup {
		log.Ctx(ctx).Warn().Msg("Google register rejected because tenant setup is not confirmed")
		return nil, codedError(
			400,
			"Tenant setup confirmation is required",
			errorCodeTenantSetupConfirmationRequired,
		)
	}

	tenantCode, err := normalizeAndValidateTenantCode(input.TenantCode)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	identity, err := c.validateGoogleIDToken(ctx, input.IDToken, input.Nonce)
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Warn().Err(err).Msg("Google register token validation failed")
		return nil, errmsg.NewCustomErrors(401).SetMessage("Invalid Google token")
	}

	var result *coreentity.GoogleRegisterResult
	err = c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		tenantExist, err := c.repo.ExistsTenantByCode(txCtx, tenantCode)
		if err != nil {
			return err
		}
		if tenantExist {
			log.Ctx(txCtx).Warn().Any(common.LogKeyPayload, map[string]string{
				"tenant_code": tenantCode,
			}).Msg("Tenant code already used")
			return codedError(400, "Tenant code is already registered", errorCodeTenantCodeAlreadyUsed)
		}

		linkedIdentity, err := c.repo.FindAuthIdentityByProviderSubject(
			txCtx,
			coreentity.AuthProviderGoogle,
			identity.Subject,
		)
		if err != nil {
			return err
		}
		if linkedIdentity != nil {
			log.Ctx(txCtx).Warn().Any(common.LogKeyPayload, map[string]string{
				"provider": coreentity.AuthProviderGoogle,
			}).Msg("Google identity already linked")
			return codedError(
				400,
				"Google identity is already linked",
				errorCodeGoogleIdentityAlreadyLinked,
			)
		}

		users, err := c.repo.FindActiveUsersByEmail(txCtx, identity.Email)
		if err != nil {
			return err
		}
		if len(users) > 0 {
			log.Ctx(txCtx).Warn().Any(common.LogKeyPayload, map[string]string{
				"email": identity.Email,
			}).Msg("Google register email already registered")
			return codedError(
				400,
				"Google email is already registered",
				errorCodeGoogleEmailAlreadyRegistered,
			)
		}

		tenantName := strings.TrimSpace(input.TenantName)
		tenantID, companyID, err := c.createTenant(txCtx, coreentity.Tenant{
			Name:              tenantName,
			Code:              tenantCode,
			PreferredLanguage: input.PreferredLanguage,
		})
		if err != nil {
			return err
		}

		ownerRole, err := c.getOwnerRole(txCtx, tenantID)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		user := coreentity.User{
			RoleIDs:           []string{ownerRole.ID},
			RoleNames:         []string{"owner"},
			Name:              displayName(identity),
			UserName:          usernameFromEmail(identity.Email),
			Email:             identity.Email,
			TenantID:          tenantID,
			TenantName:        tenantName,
			CompanyID:         &companyID,
			CompanyName:       &tenantName,
			EmailVerifiedAt:   &now,
			PreferredLanguage: input.PreferredLanguage,
		}

		userID, err := c.repo.InsertUser(txCtx, user)
		if err != nil {
			return err
		}
		user.ID = userID

		err = c.repo.CreateAuthIdentity(txCtx, coreentity.UserAuthIdentity{
			TenantID:        tenantID,
			UserID:          userID,
			Provider:        coreentity.AuthProviderGoogle,
			ProviderSubject: identity.Subject,
			Email:           identity.Email,
			EmailVerified:   identity.EmailVerified,
			DisplayName:     identity.DisplayName,
			PictureURL:      identity.PictureURL,
		})
		if err != nil {
			return err
		}

		tokens, err := c.issueAuthTokens(txCtx, user)
		if err != nil {
			return err
		}

		result = &coreentity.GoogleRegisterResult{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			TenantCode:   tenantCode,
		}
		return nil
	})
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	return result, nil
}

func (c *userCore) GoogleLogin(
	ctx context.Context,
	input coreentity.GoogleLoginInput,
) (*coreentity.AuthTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:google_auth:GoogleLogin")
	defer span.End()

	identity, err := c.validateGoogleIDToken(ctx, input.IDToken, input.Nonce)
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Warn().Err(err).Msg("Google login token validation failed")
		return nil, errmsg.NewCustomErrors(401).SetMessage("Invalid Google token")
	}

	var user *coreentity.User
	var authIdentity *coreentity.UserAuthIdentity
	err = c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		linkedIdentity, err := c.repo.FindAuthIdentityByProviderSubject(
			txCtx,
			coreentity.AuthProviderGoogle,
			identity.Subject,
		)
		if err != nil {
			return err
		}

		if linkedIdentity != nil {
			foundUser, err := c.repo.FindActiveUserByIDAndTenant(
				txCtx,
				linkedIdentity.UserID,
				linkedIdentity.TenantID,
			)
			if err != nil {
				return err
			}
			user = foundUser
			authIdentity = linkedIdentity
			return nil
		}

		users, err := c.repo.FindActiveUsersByEmail(txCtx, identity.Email)
		if err != nil {
			return err
		}

		switch len(users) {
		case 0:
			log.Ctx(txCtx).Warn().Any(common.LogKeyPayload, map[string]string{
				"email": identity.Email,
			}).Msg("Google account is not connected")
			return codedError(
				400,
				"Google account is not connected to a Presense user",
				errorCodeGoogleAccountNotConnected,
			)
		case 1:
			foundUser := users[0]
			err = c.repo.CreateAuthIdentity(txCtx, coreentity.UserAuthIdentity{
				TenantID:        foundUser.TenantID,
				UserID:          foundUser.ID,
				Provider:        coreentity.AuthProviderGoogle,
				ProviderSubject: identity.Subject,
				Email:           identity.Email,
				EmailVerified:   identity.EmailVerified,
				DisplayName:     identity.DisplayName,
				PictureURL:      identity.PictureURL,
			})
			if err != nil {
				return err
			}
			user = &foundUser
			return nil
		default:
			log.Ctx(txCtx).Warn().Any(common.LogKeyPayload, map[string]string{
				"email": identity.Email,
			}).Msg("Google email exists in multiple tenants")
			return codedError(
				400,
				"This Google email exists in multiple tenants. Sign in with email and password first.",
				errorCodeGoogleEmailAmbiguous,
			)
		}
	})
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	if authIdentity != nil {
		if err := c.repo.UpdateAuthIdentityLastLogin(ctx, authIdentity.ID); err != nil {
			tracing.RecordError(span, err)
			return nil, err
		}
	}

	return c.issueAuthTokens(ctx, *user)
}

func (c *userCore) validateGoogleIDToken(
	ctx context.Context,
	idToken string,
	nonce string,
) (*coreentity.GoogleIdentity, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:google_auth:validateGoogleIDToken")
	defer span.End()

	if c.googleTokenValidator == nil {
		return nil, errors.New("google token validator is not configured")
	}
	return c.googleTokenValidator.Validate(ctx, idToken, nonce)
}

func normalizeAndValidateTenantCode(code string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(code))
	if !tenantCodePattern.MatchString(normalized) {
		return "", codedError(400, "Tenant code is invalid", errorCodeTenantCodeInvalid)
	}
	if _, ok := tenantCodeReserved[normalized]; ok {
		return "", codedError(400, "Tenant code is reserved", errorCodeTenantCodeReserved)
	}
	return normalized, nil
}

func codedError(status int, message string, code string) *errmsg.CustomError {
	return errmsg.NewCustomErrors(status).
		SetMessage(message).
		SetErrorCode(code)
}

func displayName(identity *coreentity.GoogleIdentity) string {
	if strings.TrimSpace(identity.DisplayName) != "" {
		return strings.TrimSpace(identity.DisplayName)
	}
	return usernameFromEmail(identity.Email)
}

func usernameFromEmail(email string) string {
	local, _, found := strings.Cut(email, "@")
	if !found {
		return strings.TrimSpace(email)
	}
	local = strings.TrimSpace(local)
	if local == "" {
		return strings.TrimSpace(email)
	}
	return local
}
