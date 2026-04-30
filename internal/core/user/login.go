package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func (c *userCore) Login(ctx context.Context, user coreentity.User) (*coreentity.AuthTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:login:Login")
	defer span.End()

	foundUser, err := c.repo.FindActiveUserByEmailAndTenant(ctx, user.Email, user.TenantCode)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{"email": user.Email}).Msg("Invalid credentials")
		return nil, errmsg.NewCustomErrors(400).SetMessage("Invalid credentials")
	}

	if foundUser.EmailVerifiedAt == nil {
		log.Ctx(ctx).
			Warn().
			Any(common.LogKeyPayload, map[string]string{"email": user.Email}).
			Msg("Email is not verified")
		return nil, errmsg.NewCustomErrors(403).SetMessage("Email is not verified")
	}

	return c.issueAuthTokens(ctx, *foundUser)
}
