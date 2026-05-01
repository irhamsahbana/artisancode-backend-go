package handler

import (
	"fmt"
	"strings"
	"time"

	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
)

const (
	resendVerificationEmailCooldown = 60 * time.Second
	resendVerificationEmailLimit    = 3
	resendVerificationEmailWindow   = 15 * time.Minute
)

func (h *userHandler) limitResendVerificationEmail(c *fiber.Ctx, email, tenantCode string) error {
	if h.rateLimiter == nil {
		return nil
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	normalizedTenantCode := strings.ToUpper(strings.TrimSpace(tenantCode))
	cooldownKey := fmt.Sprintf("resend-verification-cooldown:%s:%s:%s", normalizedTenantCode, normalizedEmail, c.IP())
	allowed, retryAfter := h.rateLimiter.Allow(cooldownKey, 1, resendVerificationEmailCooldown)
	if !allowed {
		c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
		return c.Status(fiber.StatusTooManyRequests).JSON(
			response.Error(errmsg.NewCustomErrors(429).SetMessage(errmsg.MessageTooManyVerificationEmailRequestsPleaseWaitBeforeTryingAgain)),
		)
	}

	burstKey := fmt.Sprintf("resend-verification-burst:%s:%s:%s", normalizedTenantCode, normalizedEmail, c.IP())
	allowed, retryAfter = h.rateLimiter.Allow(burstKey, resendVerificationEmailLimit, resendVerificationEmailWindow)
	if allowed {
		return nil
	}

	c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
	return c.Status(fiber.StatusTooManyRequests).JSON(
		response.Error(errmsg.NewCustomErrors(429).SetMessage(errmsg.MessageTooManyVerificationEmailRequestsPleaseWaitBeforeTryingAgain)),
	)
}
