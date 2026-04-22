package core

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/config"
	"codebase-app/pkg/errmsg"
)

const (
	emailVerificationTokenTTL = time.Hour * 24 * 7
	passwordResetTokenTTL     = time.Minute * 30
)

func generateUserActionToken() (string, string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}

	rawToken := hex.EncodeToString(buf)
	hashed := sha256.Sum256([]byte(rawToken))
	return rawToken, hex.EncodeToString(hashed[:]), nil
}

func hashUserActionToken(token string) string {
	hashed := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hashed[:])
}

func (c *userCore) issueEmailVerification(ctx context.Context, user coreentity.User) error {
	rawToken, tokenHash, err := generateUserActionToken()
	if err != nil {
		return err
	}

	err = c.repo.DeleteUserActionTokensByPurpose(ctx, user.ID, coreentity.UserActionTokenPurposeEmailVerification)
	if err != nil {
		return err
	}

	err = c.repo.CreateUserActionToken(ctx, coreentity.UserActionToken{
		UserID:    user.ID,
		Purpose:   coreentity.UserActionTokenPurposeEmailVerification,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(emailVerificationTokenTTL),
	})
	if err != nil {
		return err
	}

	err = c.enqueueVerificationEmail(
		ctx,
		user.Name,
		user.Email,
		user.TenantName,
		buildFrontendActionURL(config.Envs.FrontendURL.EmailVerification, rawToken),
		user.PreferredLanguage,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *userCore) issuePasswordReset(ctx context.Context, user coreentity.User) error {
	rawToken, tokenHash, err := generateUserActionToken()
	if err != nil {
		return err
	}

	err = c.repo.DeleteUserActionTokensByPurpose(ctx, user.ID, coreentity.UserActionTokenPurposePasswordReset)
	if err != nil {
		return err
	}

	err = c.repo.CreateUserActionToken(ctx, coreentity.UserActionToken{
		UserID:    user.ID,
		Purpose:   coreentity.UserActionTokenPurposePasswordReset,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(passwordResetTokenTTL),
	})
	if err != nil {
		return err
	}

	err = c.enqueuePasswordResetEmail(
		ctx,
		user.Name,
		user.Email,
		user.TenantName,
		buildFrontendActionURL(config.Envs.FrontendURL.PasswordReset, rawToken),
		user.PreferredLanguage,
	)
	if err != nil {
		return err
	}

	return nil
}

func buildFrontendActionURL(path, rawToken string) string {
	return fmt.Sprintf("%s%s?token=%s",
		config.Envs.FrontendURL.ClientBaseURL,
		path,
		url.QueryEscape(rawToken),
	)
}

type queuedEmailPayload struct {
	Email             string `json:"email"`
	UserName          string `json:"user_name"`
	TenantName        string `json:"tenant_name"`
	ActionLink        string `json:"action_link"`
	PreferredLanguage string `json:"preferred_language"`
}

func (c *userCore) enqueueVerificationEmail(ctx context.Context, userName, email, tenantName, verificationLink, preferredLanguage string) error {
	if c.bus == nil {
		return errmsg.NewCustomErrors(500).SetMessage("Email message bus is not configured")
	}

	return c.bus.PublishJSON(ctx, common.MessageSubjectEmailVerification, queuedEmailPayload{
		Email:             email,
		UserName:          userName,
		TenantName:        tenantName,
		ActionLink:        verificationLink,
		PreferredLanguage: preferredLanguage,
	})
}

func (c *userCore) enqueuePasswordResetEmail(ctx context.Context, userName, email, tenantName, resetLink, preferredLanguage string) error {
	if c.bus == nil {
		return errmsg.NewCustomErrors(500).SetMessage("Email message bus is not configured")
	}

	return c.bus.PublishJSON(ctx, common.MessageSubjectEmailForgotPassword, queuedEmailPayload{
		Email:             email,
		UserName:          userName,
		TenantName:        tenantName,
		ActionLink:        resetLink,
		PreferredLanguage: preferredLanguage,
	})
}
