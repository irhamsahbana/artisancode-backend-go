package core

import (
	"context"
	"fmt"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/config"
	"codebase-app/internal/infrastructure/tracing"
	emailint "codebase-app/internal/integration/email"
	"codebase-app/internal/integration/tokencache"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/jwthandler"
	"codebase-app/pkg/security"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func (c *userCore) RegisterOwner(ctx context.Context, user coreentity.User, tenant coreentity.Tenant) (*coreentity.AuthTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:register:RegisterOwner")
	defer span.End()

	exist, err := c.repo.ExistsActiveUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if exist {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{"email": user.Email}).Msg("Email already registered")
		return nil, errmsg.NewCustomErrors(400).SetMessage("Email is already registered")
	}

	tenantExist, err := c.repo.ExistsTenantByCode(ctx, tenant.Code)
	if err != nil {
		return nil, err
	}
	if tenantExist {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{"tenantCode": tenant.Code}).Msg("Tenant code already registered")
		return nil, errmsg.NewCustomErrors(400).SetMessage("Tenant code is already registered")
	}

	tenantID, err := c.createTenant(ctx, tenant)
	if err != nil {
		return nil, err
	}

	ownerRole, err := c.getOwnerRole(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	userID, err := c.createOwnerUser(ctx, user, tenantID, ownerRole.ID)
	if err != nil {
		return nil, err
	}

	token, refreshToken, err := c.generateAuthTokens(ctx, userID, tenantID, user.TenantName, user.UserName, []string{ownerRole.Name})
	if err != nil {
		return nil, err
	}

	c.sendVerificationEmail(user.Name, user.Email, user.TenantName, userID, refreshToken)

	return &coreentity.AuthTokens{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (c *userCore) createTenant(ctx context.Context, tenant coreentity.Tenant) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:register:createTenant")
	defer span.End()

	tenantID, err := c.repo.InsertTenant(ctx, coreentity.Tenant{
		Name: tenant.Name,
		Code: tenant.Code,
	})
	if err != nil {
		return "", err
	}

	_, err = c.repo.InitializeTenant(ctx, tenantID, tenant.Name)
	if err != nil {
		return "", err
	}

	return tenantID, nil
}

func (c *userCore) getOwnerRole(ctx context.Context, tenantID string) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:register:getOwnerRole")
	defer span.End()

	return c.repo.GetRoleByName(ctx, "owner", tenantID)
}

func (c *userCore) createOwnerUser(ctx context.Context, user coreentity.User, tenantID, roleID string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:register:createOwnerUser")
	defer span.End()

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	userData := coreentity.User{
		RoleIDs:    []string{roleID},
		RoleNames:  []string{"owner"},
		Name:       user.Name,
		UserName:   user.UserName,
		Email:      user.Email,
		Password:   string(hashed),
		TenantID:   tenantID,
		TenantName: user.TenantName,
	}

	return c.repo.InsertUser(ctx, userData)
}

func (c *userCore) generateAuthTokens(ctx context.Context, userID, tenantID, tenantName, userName string, roles []string) (string, string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:register:generateAuthTokens")
	defer span.End()

	tokenExp := time.Now().UTC().Add(time.Hour * 24)
	payload := jwthandler.CostumClaimsPayload{
		UserID:          userID,
		TenantID:        tenantID,
		TenantName:      tenantName,
		UserName:        userName,
		Roles:           roles,
		TokenExpiration: tokenExp,
	}

	token, err := jwthandler.GenerateTokenString(payload)
	if err != nil {
		return "", "", errmsg.NewCustomErrors(500).SetMessage("Failed to generate token")
	}

	refreshToken := uuid.New().String()
	refreshTokenData := tokencache.RefreshTokenData{
		UserID:   userID,
		TenantID: tenantID,
	}
	c.tokenCache.SetRefreshToken(refreshToken, refreshTokenData, time.Hour*24*7)

	return token, refreshToken, nil
}

func (c *userCore) sendVerificationEmail(userName, email, tenantName, userID, refreshToken string) {
	verificationLinkData := fmt.Sprintf("user_id=%s&token=%s", userID, refreshToken)
	verificationURL := fmt.Sprintf("%s%s?%s",
		config.Envs.FrontendURL.ClientBaseURL,
		config.Envs.FrontendURL.EmailVerification,
		verificationLinkData,
	)

	signedURL, err := security.GenerateSignedURL(verificationURL, time.Hour*24*7)
	if err != nil {
		return
	}

	emailint.SendEmail(emailint.EmailPayload{
		To:      []string{email},
		Subject: "Email Verification - " + tenantName,
		Body:    buildVerificationEmailBody(tenantName, userName, signedURL.Link),
	})
}

func buildVerificationEmailBody(tenantName, userName, verificationLink string) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Email Verification</title>
		</head>
		<body>
			<h1>Welcome to %s</h1>
			<p>Thank you for registering, %s!</p>
			<p>Click the link below to verify your email:</p>
			<a href="%s">Verify Email</a>
			<p>This link will expire in 7 days.</p>
		</body>
		</html>
	`, tenantName, userName, verificationLink)
}
