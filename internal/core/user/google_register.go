package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

const (
	errorCodeGoogleEmailAlreadyRegistered    = "google_email_already_registered"
	errorCodeTenantSetupConfirmationRequired = "tenant_setup_confirmation_required"
)

func (c *userCore) GoogleRegister(
	ctx context.Context,
	input coreentity.GoogleRegisterInput,
) (*coreentity.GoogleRegisterResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:google_register:GoogleRegister")
	defer span.End()

	if !input.ConfirmTenantSetup {
		log.Ctx(ctx).Warn().Msg("Google register rejected because tenant setup is not confirmed")
		return nil, codedError(
			400,
			errmsg.MessageTenantSetupConfirmationIsRequired,
			errorCodeTenantSetupConfirmationRequired,
		)
	}

	tenantCode, err := normalizeAndValidateTenantCode(input.TenantCode)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	identity, err := c.resolveGoogleRegisterIdentity(ctx, input)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	result, err := c.registerGoogleTenantOwner(
		ctx,
		identity,
		googleTenantSignupInput{
			TenantName:        strings.TrimSpace(input.TenantName),
			TenantCode:        tenantCode,
			PreferredLanguage: input.PreferredLanguage,
		},
	)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	if strings.TrimSpace(input.RegistrationToken) != "" {
		c.tokenCache.DeleteGoogleRegistration(ctx, input.RegistrationToken)
	}

	return result, nil
}
