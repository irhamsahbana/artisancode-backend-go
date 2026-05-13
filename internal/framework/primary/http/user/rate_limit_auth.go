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
	authRequestCooldown = 2 * time.Second
	authRequestLimit    = 10
	authRequestWindow   = 15 * time.Minute

	registrationCooldown = 10 * time.Second
	registrationLimit    = 5
	registrationWindow   = 30 * time.Minute
)

type attemptLimitConfig struct {
	scope    string
	cooldown time.Duration
	limit    int
	window   time.Duration
	message  string
}

func (h *userHandler) limitAuthenticationRequest(c fiber.Ctx, scope string, subjectParts ...string) (bool, error) {
	return limitAttempt(
		c,
		h.rateLimiter,
		attemptLimitConfig{
			scope:    scope,
			cooldown: authRequestCooldown,
			limit:    authRequestLimit,
			window:   authRequestWindow,
			message:  errmsg.MessageTooManyAuthenticationRequestsPleaseWaitBeforeTryingAgain,
		},
		subjectParts...,
	)
}

func (h *userHandler) limitRegistrationRequest(c fiber.Ctx, scope string, subjectParts ...string) (bool, error) {
	return limitAttempt(
		c,
		h.rateLimiter,
		attemptLimitConfig{
			scope:    scope,
			cooldown: registrationCooldown,
			limit:    registrationLimit,
			window:   registrationWindow,
			message:  errmsg.MessageTooManyRegistrationRequestsPleaseWaitBeforeTryingAgain,
		},
		subjectParts...,
	)
}

func limitAttempt(
	c fiber.Ctx,
	limiter interface {
		Allow(key string, limit int, window time.Duration) (bool, time.Duration)
	},
	cfg attemptLimitConfig,
	subjectParts ...string,
) (bool, error) {
	if limiter == nil {
		return false, nil
	}

	subjectKey := buildRateLimitSubjectKey(subjectParts...)
	cooldownKey := fmt.Sprintf("%s-cooldown:%s:%s", cfg.scope, subjectKey, c.IP())
	allowed, retryAfter := limiter.Allow(cooldownKey, 1, cfg.cooldown)
	if !allowed {
		return true, tooManyRequestsResponse(c, retryAfter, cfg.message)
	}

	burstKey := fmt.Sprintf("%s-burst:%s:%s", cfg.scope, subjectKey, c.IP())
	allowed, retryAfter = limiter.Allow(burstKey, cfg.limit, cfg.window)
	if allowed {
		return false, nil
	}

	return true, tooManyRequestsResponse(c, retryAfter, cfg.message)
}

func tooManyRequestsResponse(c fiber.Ctx, retryAfter time.Duration, message string) error {
	c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
	return c.Status(fiber.StatusTooManyRequests).JSON(
		response.Error(errmsg.NewCustomErrors(429).SetMessage(message)),
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

func normalizedTenantCode(tenantCode string) string {
	return strings.ToUpper(strings.TrimSpace(tenantCode))
}
