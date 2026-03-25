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

func (c *userCore) RefreshToken(ctx context.Context, user coreentity.User) (*coreentity.AuthTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "core.RefreshToken")
	defer span.End()

	tokenData, found := c.tokenCache.GetRefreshToken(user.RefreshToken)
	if !found {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"refreshToken": user.RefreshToken,
		}).Msg("Invalid or expired refresh token")
		return nil, errmsg.NewCustomErrors(401).SetMessage("Invalid or expired refresh token")
	}

	foundUser, err := c.repo.FindActiveUserByEmailAndTenant(ctx, "", tokenData.TenantID)
	if err != nil {
		return nil, err
	}

	c.tokenCache.DeleteRefreshToken(user.RefreshToken)

	newRefreshToken := uuid.New().String()
	newRefreshTokenData := tokencache.RefreshTokenData{
		UserID:   foundUser.ID,
		TenantID: foundUser.TenantID,
	}
	c.tokenCache.SetRefreshToken(newRefreshToken, newRefreshTokenData, time.Hour*24*7)

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

	newToken, err := jwthandler.GenerateTokenString(payload)
	if err != nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage("Failed to generate token")
	}

	return &coreentity.AuthTokens{
		AccessToken:  newToken,
		RefreshToken: newRefreshToken,
	}, nil
}
