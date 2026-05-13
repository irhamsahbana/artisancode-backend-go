package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
)

type stubAttemptLimiter struct {
	responses  []bool
	retryAfter time.Duration
	calls      int
}

func (s *stubAttemptLimiter) Allow(_ string, _ int, _ time.Duration) (bool, time.Duration) {
	if s.calls >= len(s.responses) {
		return true, 0
	}

	allowed := s.responses[s.calls]
	s.calls++
	if allowed {
		return true, 0
	}

	return false, s.retryAfter
}

func TestLimitAuthenticationRequestReturnsTooManyRequests(t *testing.T) {
	h := &internalUserHandler{
		rateLimiter: &stubAttemptLimiter{
			responses:  []bool{false},
			retryAfter: 12 * time.Second,
		},
	}

	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		if handled, err := h.limitAuthenticationRequest(c, "internal-user-login", "admin@example.com"); err != nil || handled {
			return err
		}
		return c.SendStatus(fiber.StatusNoContent)
	})

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != fiber.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusTooManyRequests)
	}

	if got := resp.Header.Get("Retry-After"); got != "12" {
		t.Fatalf("Retry-After = %q, want %q", got, "12")
	}
}
