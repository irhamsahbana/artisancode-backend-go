package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

const (
	errorCodeGoogleAccountNotConnected   = "google_account_not_connected"
	errorCodeGoogleEmailAmbiguous        = "google_email_ambiguous"
	errorCodeGoogleIdentityAlreadyLinked = "google_identity_already_linked"
)

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
		return nil, errmsg.NewCustomErrors(401).SetMessage(errmsg.MessageInvalidGoogleToken)
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
				errmsg.MessageGoogleAccountIsNotConnectedToAPresenseUser,
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
				errmsg.MessageThisGoogleEmailExistsInMultipleTenantsSignInWithEmailAndPasswordFirst,
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
