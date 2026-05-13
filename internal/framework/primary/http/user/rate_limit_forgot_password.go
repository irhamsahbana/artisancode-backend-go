package handler

import (
	"fmt"
	"strings"
	"time"

	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
)

const (
	forgotPasswordCooldown = 60 * time.Second
	forgotPasswordLimit    = 3
	forgotPasswordWindow   = 15 * time.Minute
)

func (h *userHandler) limitForgotPassword(c fiber.Ctx, email, tenantCode string) error {
	if h.rateLimiter == nil {
		return nil
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	normalizedTenantCode := strings.ToUpper(strings.TrimSpace(tenantCode))
	cooldownKey := fmt.Sprintf("forgot-password-cooldown:%s:%s:%s", normalizedTenantCode, normalizedEmail, c.IP())
	allowed, retryAfter := h.rateLimiter.Allow(cooldownKey, 1, forgotPasswordCooldown)
	if !allowed {
		c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
		return c.Status(fiber.StatusTooManyRequests).JSON(
			response.Error(errmsg.NewCustomErrors(429).SetMessage(errmsg.MessageTooManyPasswordResetRequestsPleaseWaitBeforeTryingAgain)),
		)
	}

	burstKey := fmt.Sprintf("forgot-password-burst:%s:%s:%s", normalizedTenantCode, normalizedEmail, c.IP())
	allowed, retryAfter = h.rateLimiter.Allow(burstKey, forgotPasswordLimit, forgotPasswordWindow)
	if allowed {
		return nil
	}

	c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
	return c.Status(fiber.StatusTooManyRequests).JSON(
		response.Error(errmsg.NewCustomErrors(429).SetMessage(errmsg.MessageTooManyPasswordResetRequestsPleaseWaitBeforeTryingAgain)),
	)
}
