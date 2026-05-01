package core

import (
	"context"
	"strings"
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

func (c *internalUserCore) Login(ctx context.Context, user coreentity.InternalUser) (*coreentity.AuthTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:core:Login")
	defer span.End()

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	foundUser, err := c.repo.FindActiveInternalUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if foundUser.Status != coreentity.InternalUserStatusActive {
		return nil, errmsg.NewCustomErrors(403).SetMessage(errmsg.MessageInternalUserIsNotActive)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		log.Ctx(ctx).
			Warn().
			Any(common.LogKeyPayload, map[string]string{"email": user.Email}).
			Msg("Invalid internal credentials")
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInvalidCredentials)
	}

	tokenExp := time.Now().UTC().Add(15 * time.Minute)
	token, err := jwthandler.GenerateInternalUserTokenString(ctx, jwthandler.InternalUserClaimsPayload{
		UserID:          foundUser.ID,
		UserName:        foundUser.FullName,
		Roles:           []string{foundUser.RoleCode},
		TokenExpiration: tokenExp,
	})
	if err != nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageFailedToGenerateToken)
	}

	refreshToken := uuid.New().String()
	c.tokenCache.SetRefreshToken(ctx, refreshToken, tokencache.RefreshTokenData{
		UserID: foundUser.ID,
	}, 7*24*time.Hour)

	if err := c.repo.UpdateInternalUserLastLogin(ctx, foundUser.ID); err != nil {
		return nil, err
	}

	return &coreentity.AuthTokens{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}
