package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/internal/integration/tokencache"
	"codebase-app/pkg/errmsg"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const googleRegistrationTokenTTL = 10 * time.Minute

func (c *userCore) GoogleRegisterInit(
	ctx context.Context,
	input coreentity.GoogleRegisterInitInput,
) (*coreentity.GoogleRegisterInitResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:google_register_init:GoogleRegisterInit")
	defer span.End()

	identity, err := c.validateGoogleIDToken(ctx, input.IDToken, input.Nonce)
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Warn().Err(err).Msg("Google register init token validation failed")
		return nil, errmsg.NewCustomErrors(401).SetMessage(errmsg.MessageInvalidGoogleToken)
	}

	registrationToken := uuid.NewString()
	c.tokenCache.SetGoogleRegistration(ctx, registrationToken, tokencache.GoogleRegistrationData{
		Subject:       identity.Subject,
		Email:         identity.Email,
		EmailVerified: identity.EmailVerified,
		DisplayName:   identity.DisplayName,
		PictureURL:    identity.PictureURL,
		Nonce:         identity.Nonce,
	}, googleRegistrationTokenTTL)

	return &coreentity.GoogleRegisterInitResult{
		RegistrationToken: registrationToken,
		Email:             identity.Email,
		DisplayName:       identity.DisplayName,
		PictureURL:        identity.PictureURL,
	}, nil
}
