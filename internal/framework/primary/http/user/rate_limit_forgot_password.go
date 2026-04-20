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
	forgotPasswordCooldown = 60 * time.Second
	forgotPasswordLimit    = 3
	forgotPasswordWindow   = 15 * time.Minute
)

func (h *userHandler) limitForgotPassword(c *fiber.Ctx, email string) error {
	if h.rateLimiter == nil {
		return nil
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	cooldownKey := fmt.Sprintf("forgot-password-cooldown:%s:%s", normalizedEmail, c.IP())
	allowed, retryAfter := h.rateLimiter.Allow(cooldownKey, 1, forgotPasswordCooldown)
	if !allowed {
		c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
		return c.Status(fiber.StatusTooManyRequests).JSON(
			response.Error(errmsg.NewCustomErrors(429).SetMessage("Too many password reset requests. Please wait before trying again")),
		)
	}

	burstKey := fmt.Sprintf("forgot-password-burst:%s:%s", normalizedEmail, c.IP())
	allowed, retryAfter = h.rateLimiter.Allow(burstKey, forgotPasswordLimit, forgotPasswordWindow)
	if allowed {
		return nil
	}

	c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
	return c.Status(fiber.StatusTooManyRequests).JSON(
		response.Error(errmsg.NewCustomErrors(429).SetMessage("Too many password reset requests. Please wait before trying again")),
	)
}
