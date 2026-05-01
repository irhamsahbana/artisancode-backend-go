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
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:refresh_token:RefreshToken")
	defer span.End()

	tokenData, found := c.tokenCache.GetRefreshToken(ctx, user.RefreshToken)
	if !found {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"refreshToken": user.RefreshToken,
		}).Msg(errmsg.MessageInvalidOrExpiredRefreshToken)
		return nil, errmsg.NewCustomErrors(401).SetMessage(errmsg.MessageInvalidOrExpiredRefreshToken)
	}

	foundUser, err := c.repo.FindActiveUserByIDAndTenant(ctx, tokenData.UserID, tokenData.TenantID)
	if err != nil {
		return nil, err
	}

	c.tokenCache.DeleteRefreshToken(ctx, user.RefreshToken)

	newRefreshToken := uuid.New().String()
	newRefreshTokenData := tokencache.RefreshTokenData{
		UserID:   foundUser.ID,
		TenantID: foundUser.TenantID,
	}
	c.tokenCache.SetRefreshToken(ctx, newRefreshToken, newRefreshTokenData, time.Hour*24*7)

	tokenExp := time.Now().UTC().Add(time.Minute * 15)
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

	newToken, err := jwthandler.GenerateTokenString(ctx, payload)
	if err != nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageFailedToGenerateToken)
	}

	return &coreentity.AuthTokens{
		AccessToken:  newToken,
		RefreshToken: newRefreshToken,
	}, nil
}
