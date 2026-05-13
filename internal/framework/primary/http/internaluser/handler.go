package handler

import (
	"codebase-app/internal/integration/ratelimit"
	"codebase-app/internal/middleware"
	corePorts "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v3"
)

type internalUserHandler struct {
	core        corePorts.InternalUserCore
	rateLimiter ratelimit.AttemptLimiter
}

type Config struct {
	Core        corePorts.InternalUserCore
	RateLimiter ratelimit.AttemptLimiter
}

func NewInternalUserHandler(cfg Config) *internalUserHandler {
	return &internalUserHandler{
		core:        cfg.Core,
		rateLimiter: cfg.RateLimiter,
	}
}

func (h *internalUserHandler) Register(router fiber.Router) {
	router.Post("/login", h.login)
	router.Post("/refresh-token", h.refreshToken)
	router.Post("/logout", h.logout)

	router.Get("/", middleware.InternalAuth, h.getInternalUsers)
	router.Get("/:id", middleware.InternalAuth, h.getInternalUser)
	router.Post("/", middleware.InternalAuth, h.createInternalUser)
	router.Put("/:id", middleware.InternalAuth, h.updateInternalUser)
	router.Delete("/:id", middleware.InternalAuth, h.deleteInternalUser)
}
