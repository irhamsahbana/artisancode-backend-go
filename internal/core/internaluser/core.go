package core

import (
	"context"
	"strings"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/internal/integration/tokencache"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/jwthandler"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type internalUserCore struct {
	repo       portsRepo.InternalUserRepository
	tokenCache tokencache.TokenCacheContract
}

type Config struct {
	Repo       portsRepo.InternalUserRepository
	TokenCache tokencache.TokenCacheContract
}

var _ corePorts.InternalUserCore = &internalUserCore{}

func NewInternalUserCore(cfg Config) *internalUserCore {
	return &internalUserCore{
		repo:       cfg.Repo,
		tokenCache: cfg.TokenCache,
	}
}

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

	tokenExp := time.Now().UTC().Add(24 * time.Hour)
	token, err := jwthandler.GenerateInternalUserTokenString(jwthandler.InternalUserClaimsPayload{
		UserID:          foundUser.ID,
		UserName:        foundUser.FullName,
		Roles:           []string{foundUser.RoleCode},
		TokenExpiration: tokenExp,
	})
	if err != nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageFailedToGenerateToken)
	}

	refreshToken := uuid.New().String()
	c.tokenCache.SetRefreshToken(refreshToken, tokencache.RefreshTokenData{
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

func (c *internalUserCore) RefreshToken(
	ctx context.Context,
	user coreentity.InternalUser,
) (*coreentity.AuthTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:core:RefreshToken")
	defer span.End()

	tokenData, found := c.tokenCache.GetRefreshToken(user.RefreshToken)
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

	token, err := jwthandler.GenerateInternalUserTokenString(jwthandler.InternalUserClaimsPayload{
		UserID:          foundUser.ID,
		UserName:        foundUser.FullName,
		Roles:           []string{foundUser.RoleCode},
		TokenExpiration: time.Now().UTC().Add(24 * time.Hour),
	})
	if err != nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageFailedToGenerateToken)
	}

	newRefreshToken := uuid.New().String()
	c.tokenCache.DeleteRefreshToken(user.RefreshToken)
	c.tokenCache.SetRefreshToken(newRefreshToken, tokencache.RefreshTokenData{
		UserID: foundUser.ID,
	}, 7*24*time.Hour)

	return &coreentity.AuthTokens{
		AccessToken:  token,
		RefreshToken: newRefreshToken,
	}, nil
}

func (c *internalUserCore) GetInternalUsers(
	ctx context.Context,
	filter coreentity.InternalUserListFilter,
) ([]coreentity.InternalUser, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:core:GetInternalUsers")
	defer span.End()

	if err := c.authorizeManage(common.GetUserContext(ctx)); err != nil {
		return nil, 0, err
	}

	return c.repo.GetInternalUsers(ctx, filter)
}

func (c *internalUserCore) GetInternalUser(
	ctx context.Context,
	filter coreentity.InternalUserFilter,
) (*coreentity.InternalUser, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:core:GetInternalUser")
	defer span.End()

	if err := c.authorizeManage(common.GetUserContext(ctx)); err != nil {
		return nil, err
	}

	return c.repo.GetInternalUser(ctx, filter)
}

func (c *internalUserCore) CreateInternalUser(
	ctx context.Context,
	data coreentity.InternalUser,
) (*coreentity.InternalUser, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:core:CreateInternalUser")
	defer span.End()

	if err := c.authorizeManage(data.UserCtx); err != nil {
		return nil, err
	}

	if err := c.normalizeAndValidate(&data, true); err != nil {
		return nil, err
	}

	exists, err := c.repo.ExistsInternalUserByEmail(ctx, data.Email, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalUserEmailAlreadyExists)
	}

	hashedPassword, err := hashInternalPassword(data.Password)
	if err != nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageFailedToHashPassword)
	}
	data.Password = hashedPassword

	return c.repo.CreateInternalUser(ctx, data)
}

func (c *internalUserCore) UpdateInternalUser(ctx context.Context, data coreentity.InternalUser) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:core:UpdateInternalUser")
	defer span.End()

	if err := c.authorizeManage(data.UserCtx); err != nil {
		return err
	}

	if err := c.normalizeAndValidate(&data, false); err != nil {
		return err
	}

	exists, err := c.repo.ExistsInternalUserByEmail(ctx, data.Email, data.ID)
	if err != nil {
		return err
	}
	if exists {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalUserEmailAlreadyExists)
	}

	if data.Password != "" {
		hashedPassword, err := hashInternalPassword(data.Password)
		if err != nil {
			return errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageFailedToHashPassword)
		}
		data.Password = hashedPassword
	}

	return c.repo.UpdateInternalUser(ctx, data)
}

func (c *internalUserCore) DeleteInternalUser(ctx context.Context, filter coreentity.InternalUserDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:core:DeleteInternalUser")
	defer span.End()

	if err := c.authorizeManage(common.GetUserContext(ctx)); err != nil {
		return err
	}

	return c.repo.DeleteInternalUser(ctx, filter)
}

func (c *internalUserCore) authorizeManage(userCtx common.UserContext) error {
	if !userCtx.HasRole(coreentity.InternalUserRoleSuperAdmin) {
		return errmsg.NewCustomErrors(403).SetMessage(errmsg.MessageYouAreNotAuthorizedToManageInternalUsers)
	}
	return nil
}

func (c *internalUserCore) normalizeAndValidate(data *coreentity.InternalUser, passwordRequired bool) error {
	data.FullName = strings.TrimSpace(data.FullName)
	data.Email = strings.ToLower(strings.TrimSpace(data.Email))
	data.RoleCode = strings.TrimSpace(data.RoleCode)
	data.Status = strings.TrimSpace(data.Status)

	if data.RoleCode != coreentity.InternalUserRoleSuperAdmin && data.RoleCode != coreentity.InternalUserRoleOperator {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalUserRoleIsInvalid)
	}
	if data.Status != coreentity.InternalUserStatusInvited &&
		data.Status != coreentity.InternalUserStatusActive &&
		data.Status != coreentity.InternalUserStatusInactive {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalUserStatusIsInvalid)
	}
	if passwordRequired && strings.TrimSpace(data.Password) == "" {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessagePasswordIsRequired)
	}

	return nil
}

func hashInternalPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
