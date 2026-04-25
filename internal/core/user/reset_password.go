package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"golang.org/x/crypto/bcrypt"
)

func (c *userCore) ResetPassword(ctx context.Context, token coreentity.UserActionToken, user coreentity.User) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:reset_password:ResetPassword")
	defer span.End()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var resetUserID string
	err = c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		foundToken, err := c.repo.GetValidUserActionToken(txCtx, hashUserActionToken(token.Token), coreentity.UserActionTokenPurposePasswordReset)
		if err != nil {
			return err
		}
		resetUserID = foundToken.UserID

		err = c.repo.UpdateUserPassword(txCtx, foundToken.UserID, foundToken.TenantID, string(hashedPassword))
		if err != nil {
			return err
		}

		return c.repo.MarkUserActionTokenUsed(txCtx, foundToken.ID)
	})
	if err != nil {
		return err
	}

	if resetUserID != "" {
		c.tokenCache.DeleteUserRefreshTokens(resetUserID)
	}
	return nil
}
