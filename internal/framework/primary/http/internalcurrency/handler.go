package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type internalCurrencyHandler struct {
	core portsCore.InternalCurrencyCore
}

type Config struct {
	Core portsCore.InternalCurrencyCore
}

func NewInternalCurrencyHandler(cfg Config) *internalCurrencyHandler {
	return &internalCurrencyHandler{core: cfg.Core}
}

func (h *internalCurrencyHandler) Register(router fiber.Router) {
	router.Get("/providers/:provider", h.getProviderCurrencies)
	router.Put("/providers/:provider/:currency_code", h.upsertProviderCurrency)
	router.Delete("/providers/:provider/:currency_code", h.deleteProviderCurrency)

	router.Get("/", h.getInternalCurrencies)
	router.Get("/:code", h.getInternalCurrency)
	router.Post("/", h.createInternalCurrency)
	router.Put("/:code", h.updateInternalCurrency)
	router.Delete("/:code", h.deleteInternalCurrency)
}
