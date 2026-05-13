package handler

import (
	"codebase-app/internal/module/z_template_v2/ports"

	"github.com/gofiber/fiber/v3"
)

type xxxHandler struct {
	service ports.XxxService
}

func NewXxxHandler(service ports.XxxService) *xxxHandler {
	return &xxxHandler{
		service: service,
	}
}

func (h *xxxHandler) Register(router fiber.Router) {

}
