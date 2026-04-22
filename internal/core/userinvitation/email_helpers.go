package core

import (
	"context"
	"fmt"
	"net/url"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/config"

	"github.com/rs/zerolog/log"
)

type queuedInvitationEmailPayload struct {
	Email             string `json:"email"`
	UserName          string `json:"user_name"`
	TenantName        string `json:"tenant_name"`
	ActionLink        string `json:"action_link"`
	PreferredLanguage string `json:"preferred_language"`
}

func buildInvitationActionURL(rawToken string) string {
	return fmt.Sprintf("%s%s?token=%s",
		config.Envs.FrontendURL.ClientBaseURL,
		config.Envs.FrontendURL.Invitation,
		url.QueryEscape(rawToken),
	)
}

func (c *userInvitationCore) trySendInvitationEmail(ctx context.Context, item *coreentity.UserInvitation) bool {
	if item == nil {
		return false
	}

	if c.bus == nil {
		log.Ctx(ctx).Warn().Str("invitation_id", item.ID).Msg("Invitation email bus is not configured")
		return false
	}

	preferredLanguage, err := c.userRepo.GetTenantPreferredLanguage(ctx, item.TenantID)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Str("tenant_id", item.TenantID).Msg("Failed to get tenant preferred language for invitation email")
		preferredLanguage = ""
	}

	userName := item.Email
	if item.EmployeeName != nil && *item.EmployeeName != "" {
		userName = *item.EmployeeName
	}

	err = c.bus.PublishJSON(ctx, common.MessageSubjectEmailInvitation, queuedInvitationEmailPayload{
		Email:             item.Email,
		UserName:          userName,
		TenantName:        item.TenantName,
		ActionLink:        buildInvitationActionURL(item.AcceptToken),
		PreferredLanguage: preferredLanguage,
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("invitation_id", item.ID).Msg("Failed to enqueue invitation email")
		return false
	}

	return true
}
