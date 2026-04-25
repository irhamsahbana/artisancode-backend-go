package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *userCore) VerifyEmail(ctx context.Context, token coreentity.UserActionToken) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:verify_email:VerifyEmail")
	defer span.End()

	return c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		foundToken, err := c.repo.GetValidUserActionToken(txCtx, hashUserActionToken(token.Token), coreentity.UserActionTokenPurposeEmailVerification)
		if err != nil {
			return err
		}

		err = c.repo.MarkUserEmailVerified(txCtx, foundToken.UserID)
		if err != nil {
			return err
		}

		return c.repo.MarkUserActionTokenUsed(txCtx, foundToken.ID)
	})
}
