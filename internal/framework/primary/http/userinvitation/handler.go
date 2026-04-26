package handler

import (
	"codebase-app/internal/middleware"
	corePorts "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type userInvitationHandler struct {
	core corePorts.UserInvitationCore
}

type Config struct {
	Core corePorts.UserInvitationCore
}

func NewUserInvitationHandler(cfg Config) *userInvitationHandler {
	return &userInvitationHandler{core: cfg.Core}
}

func (h *userInvitationHandler) Register(router fiber.Router) {
	router.Get("/accept", h.getInvitationPreview)
	router.Post("/accept", h.acceptInvitation)
	router.Get("/", middleware.Auth, h.getInvitations)
	router.Post("/", middleware.Auth, h.createInvitation)
	router.Post("/:id/resend", middleware.Auth, h.resendInvitation)
	router.Post("/:id/revoke", middleware.Auth, h.revokeInvitation)
}
