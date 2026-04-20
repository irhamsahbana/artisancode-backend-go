package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *userCore) ResendVerificationEmail(ctx context.Context, user coreentity.User) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:resend_verification_email:ResendVerificationEmail")
	defer span.End()

	foundUser, err := c.repo.FindActiveUserByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	if foundUser.EmailVerifiedAt != nil {
		return errmsg.NewCustomErrors(400).SetMessage("Email is already verified")
	}

	foundUser.PreferredLanguage, err = c.repo.GetTenantPreferredLanguage(ctx, foundUser.TenantID)
	if err != nil {
		return err
	}

	return c.issueEmailVerification(ctx, *foundUser)
}
