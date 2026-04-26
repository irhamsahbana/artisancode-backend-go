package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/internal/integration/tokencache"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/jwthandler"

	"github.com/google/uuid"
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

	tokenExp := time.Now().UTC().Add(time.Hour * 24)
	payload := jwthandler.CostumClaimsPayload{
		UserID:          foundUser.ID,
		TenantID:        foundUser.TenantID,
		TenantName:      foundUser.TenantName,
		UserName:        foundUser.UserName,
		Roles:           foundUser.RoleNames,
		CompanyID:       foundUser.CompanyID,
		CompanyName:     foundUser.CompanyName,
		TokenExpiration: tokenExp,
	}

	token, err := jwthandler.GenerateTokenString(payload)
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"email": user.Email}).
			Msg("Failed to generate login token")
		return nil, errmsg.NewCustomErrors(500).SetMessage("Failed to generate token")
	}

	refreshToken := uuid.New().String()
	refreshTokenData := tokencache.RefreshTokenData{
		UserID:   foundUser.ID,
		TenantID: foundUser.TenantID,
	}
	c.tokenCache.SetRefreshToken(refreshToken, refreshTokenData, time.Hour*24*7)

	return &coreentity.AuthTokens{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}
