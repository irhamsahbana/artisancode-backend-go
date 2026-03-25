package handler

import (
	"codebase-app/internal/module/z_template_v2/core"
	"codebase-app/internal/module/z_template_v2/ports"
	"codebase-app/internal/module/z_template_v2/repository"

	"github.com/gofiber/fiber/v2"
)

type xxxHandler struct {
	service ports.XxxService
}

func NewXxxHandler() *xxxHandler {
	var (
		repo    = repository.NewXxxRepository()
		svc     = core.NewXxxCore(repo)
		handler = new(xxxHandler)
	)
	handler.service = svc

	return handler
}

func (h *xxxHandler) Register(router fiber.Router) {

}