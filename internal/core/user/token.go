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
)

func (c *userCore) issueAuthTokens(ctx context.Context, user coreentity.User) (*coreentity.AuthTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:token:issueAuthTokens")
	defer span.End()

	tokenExp := time.Now().UTC().Add(time.Hour * 24)
	payload := jwthandler.CostumClaimsPayload{
		UserID:          user.ID,
		TenantID:        user.TenantID,
		TenantName:      user.TenantName,
		UserName:        user.UserName,
		Roles:           user.RoleNames,
		CompanyID:       user.CompanyID,
		CompanyName:     user.CompanyName,
		TokenExpiration: tokenExp,
	}

	token, err := jwthandler.GenerateTokenString(payload)
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{
				"user_id":   user.ID,
				"tenant_id": user.TenantID,
			}).
			Msg("Failed to generate auth token")
		return nil, errmsg.NewCustomErrors(500).SetMessage("Failed to generate token")
	}

	refreshToken := uuid.New().String()
	refreshTokenData := tokencache.RefreshTokenData{
		UserID:   user.ID,
		TenantID: user.TenantID,
	}
	c.tokenCache.SetRefreshToken(refreshToken, refreshTokenData, time.Hour*24*7)

	return &coreentity.AuthTokens{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}
