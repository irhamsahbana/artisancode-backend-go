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
	internalAuthRequestCooldown = 2 * time.Second
	internalAuthRequestLimit    = 10
	internalAuthRequestWindow   = 15 * time.Minute
)

func (h *internalUserHandler) limitAuthenticationRequest(c fiber.Ctx, scope string, subjectParts ...string) (bool, error) {
	if h.rateLimiter == nil {
		return false, nil
	}

	subjectKey := buildRateLimitSubjectKey(subjectParts...)
	cooldownKey := fmt.Sprintf("%s-cooldown:%s:%s", scope, subjectKey, c.IP())
	allowed, retryAfter := h.rateLimiter.Allow(cooldownKey, 1, internalAuthRequestCooldown)
	if !allowed {
		c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
		return true, c.Status(fiber.StatusTooManyRequests).JSON(
			response.Error(errmsg.NewCustomErrors(429).SetMessage(errmsg.MessageTooManyAuthenticationRequestsPleaseWaitBeforeTryingAgain)),
		)
	}

	burstKey := fmt.Sprintf("%s-burst:%s:%s", scope, subjectKey, c.IP())
	allowed, retryAfter = h.rateLimiter.Allow(burstKey, internalAuthRequestLimit, internalAuthRequestWindow)
	if allowed {
		return false, nil
	}

	c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
	return true, c.Status(fiber.StatusTooManyRequests).JSON(
		response.Error(errmsg.NewCustomErrors(429).SetMessage(errmsg.MessageTooManyAuthenticationRequestsPleaseWaitBeforeTryingAgain)),
	)
}

func buildRateLimitSubjectKey(parts ...string) string {
	normalized := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}

	if len(normalized) == 0 {
		return "anonymous"
	}

	return strings.Join(normalized, ":")
}

func normalizedEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
