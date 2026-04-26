package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type companyHandler struct {
	core portsCore.CompanyCore
}

type Config struct {
	Core portsCore.CompanyCore
}

func NewCompanyHandler(cfg Config) *companyHandler {
	return &companyHandler{core: cfg.Core}
}

func (h *companyHandler) Register(router fiber.Router) {
	router.Get("/", h.getCompanies)
	router.Get("/:id", h.getCompany)
	router.Post("/", h.createCompany)
	router.Put("/:id", h.updateCompany)
	router.Delete("/:id", h.deleteCompany)
}
