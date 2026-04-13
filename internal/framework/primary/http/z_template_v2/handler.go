package handler

import (
	core "codebase-app/internal/core/z_template_v2"
	repository "codebase-app/internal/framework/secondary/db/postgres/z_template_v2"
	"codebase-app/internal/module/z_template_v2/ports"

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
