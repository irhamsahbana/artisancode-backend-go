package core

import (
	"context"
	"errors"
	"strings"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

const errorCodeGoogleRegistrationSessionInvalid = "google_registration_session_invalid"

func (c *userCore) resolveGoogleRegisterIdentity(
	ctx context.Context,
	input coreentity.GoogleRegisterInput,
) (*coreentity.GoogleIdentity, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:google_identity:resolveGoogleRegisterIdentity")
	defer span.End()

	if strings.TrimSpace(input.RegistrationToken) != "" {
		cachedIdentity, found := c.tokenCache.GetGoogleRegistration(input.RegistrationToken)
		if !found {
			log.Ctx(ctx).Warn().Msg("Google registration session is missing or expired")
			return nil, codedError(
				400,
				errmsg.MessageGoogleRegistrationSessionIsInvalidOrExpired,
				errorCodeGoogleRegistrationSessionInvalid,
			)
		}

		return &coreentity.GoogleIdentity{
			Subject:       cachedIdentity.Subject,
			Email:         cachedIdentity.Email,
			EmailVerified: cachedIdentity.EmailVerified,
			DisplayName:   cachedIdentity.DisplayName,
			PictureURL:    cachedIdentity.PictureURL,
			Nonce:         cachedIdentity.Nonce,
		}, nil
	}

	identity, err := c.validateGoogleIDToken(ctx, input.IDToken, input.Nonce)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Google register token validation failed")
		return nil, errmsg.NewCustomErrors(401).SetMessage(errmsg.MessageInvalidGoogleToken)
	}

	return identity, nil
}

func (c *userCore) validateGoogleIDToken(
	ctx context.Context,
	idToken string,
	nonce string,
) (*coreentity.GoogleIdentity, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:google_identity:validateGoogleIDToken")
	defer span.End()

	if c.googleTokenValidator == nil {
		return nil, errors.New("google token validator is not configured")
	}
	return c.googleTokenValidator.Validate(ctx, idToken, nonce)
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
