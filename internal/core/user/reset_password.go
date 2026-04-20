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

	foundToken, err := c.repo.GetValidUserActionToken(ctx, hashUserActionToken(token.Token), coreentity.UserActionTokenPurposePasswordReset)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = c.repo.UpdateUserPassword(ctx, foundToken.UserID, foundToken.TenantID, string(hashedPassword))
	if err != nil {
		return err
	}

	err = c.repo.MarkUserActionTokenUsed(ctx, foundToken.ID)
	if err != nil {
		return err
	}

	c.tokenCache.DeleteUserRefreshTokens(foundToken.UserID)
	return nil
}
