package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/internal/integration/tokencache"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/jwthandler"

	"github.com/google/uuid"
)

func (c *internalUserCore) RefreshToken(
	ctx context.Context,
	user coreentity.InternalUser,
) (*coreentity.AuthTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:core:RefreshToken")
	defer span.End()

	tokenData, found := c.tokenCache.GetRefreshToken(ctx, user.RefreshToken)
	if !found {
		return nil, errmsg.NewCustomErrors(401).SetMessage(errmsg.MessageInvalidRefreshToken)
	}

	foundUser, err := c.repo.GetInternalUser(ctx, coreentity.InternalUserFilter{ID: tokenData.UserID})
	if err != nil {
		return nil, err
	}

	if foundUser.Status != coreentity.InternalUserStatusActive {
		return nil, errmsg.NewCustomErrors(403).SetMessage(errmsg.MessageInternalUserIsNotActive)
	}

	token, err := jwthandler.GenerateInternalUserTokenString(ctx, jwthandler.InternalUserClaimsPayload{
		UserID:          foundUser.ID,
		UserName:        foundUser.FullName,
		Roles:           []string{foundUser.RoleCode},
		TokenExpiration: time.Now().UTC().Add(15 * time.Minute),
	})
	if err != nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageFailedToGenerateToken)
	}

	newRefreshToken := uuid.New().String()
	c.tokenCache.DeleteRefreshToken(ctx, user.RefreshToken)
	c.tokenCache.SetRefreshToken(ctx, newRefreshToken, tokencache.RefreshTokenData{
		UserID: foundUser.ID,
	}, 7*24*time.Hour)

	return &coreentity.AuthTokens{
		AccessToken:  token,
		RefreshToken: newRefreshToken,
	}, nil
}
