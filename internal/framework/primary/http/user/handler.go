package handler

import (
	"codebase-app/internal/integration/ratelimit"
	"codebase-app/internal/middleware"
	corePorts "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v3"
)

type userHandler struct {
	core        corePorts.UserCore
	rateLimiter ratelimit.AttemptLimiter
}

type Config struct {
	Core        corePorts.UserCore
	RateLimiter ratelimit.AttemptLimiter
}

func NewUserHandler(cfg Config) *userHandler {
	return &userHandler{
		core:        cfg.Core,
		rateLimiter: cfg.RateLimiter,
	}
}

func (h *userHandler) Register(router fiber.Router) {
	router.Post("/login", h.login)
	router.Post("/register", h.register)
	router.Post("/google/register/init", h.googleRegisterInit)
	router.Post("/google/register", h.googleRegister)
	router.Post("/google/login", h.googleLogin)
	router.Post("/verify-email", h.verifyEmail)
	router.Post("/resend-verification-email", h.resendVerificationEmail)
	router.Post("/forgot-password", h.forgotPassword)
	router.Post("/reset-password", h.resetPassword)
	router.Post("/refresh-token", h.refreshToken)
	router.Post("/logout", h.logout)
	router.Get("/", middleware.Auth, h.getUsers)
	router.Get("/:id", middleware.Auth, h.getUser)
	router.Post("/", middleware.Auth, h.createUser)
	router.Put("/:id", middleware.Auth, h.updateUser)
	router.Delete("/:id", middleware.Auth, h.deleteUser)
}

func (h *userHandler) RegisterTenant(router fiber.Router) {
	router.Get("/profile", h.getTenantProfile)
}
