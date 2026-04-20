package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *userCore) ForgotPassword(ctx context.Context, user coreentity.User) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:forgot_password:ForgotPassword")
	defer span.End()

	foundUser, err := c.repo.FindActiveUserByEmail(ctx, user.Email)
	if err != nil {
		if customErr, ok := err.(*errmsg.CustomError); ok && customErr.Msg == "User not found" {
			return nil
		}
		return err
	}

	foundUser.PreferredLanguage, err = c.repo.GetTenantPreferredLanguage(ctx, foundUser.TenantID)
	if err != nil {
		return err
	}

	return c.issuePasswordReset(ctx, *foundUser)
}
