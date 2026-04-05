package handler

import (
	"codebase-app/internal/middleware"
	corePorts "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	core corePorts.UserCore
}

type UserHandlerConfig struct {
	Core corePorts.UserCore
}

func NewUserHandler(cfg UserHandlerConfig) *userHandler {
	return &userHandler{
		core: cfg.Core,
	}
}

func (h *userHandler) Register(router fiber.Router) {
	router.Post("/login", h.login)
	router.Post("/register", h.register)
	router.Post("/refresh-token", h.refreshToken)
	router.Get("/", middleware.Auth, h.getUsers)
	router.Get("/:id", middleware.Auth, h.getUser)
	router.Post("/", middleware.Auth, h.createUser)
	router.Put("/:id", middleware.Auth, h.updateUser)
	router.Delete("/:id", middleware.Auth, h.deleteUser)
}
