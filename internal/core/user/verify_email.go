package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *userCore) VerifyEmail(ctx context.Context, token coreentity.UserActionToken) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:verify_email:VerifyEmail")
	defer span.End()

	foundToken, err := c.repo.GetValidUserActionToken(ctx, hashUserActionToken(token.Token), coreentity.UserActionTokenPurposeEmailVerification)
	if err != nil {
		return err
	}

	err = c.repo.MarkUserEmailVerified(ctx, foundToken.UserID)
	if err != nil {
		return err
	}

	return c.repo.MarkUserActionTokenUsed(ctx, foundToken.ID)
}
