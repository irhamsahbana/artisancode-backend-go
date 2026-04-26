package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type internalProductHandler struct {
	core portsCore.InternalProductCore
}

type Config struct {
	Core portsCore.InternalProductCore
}

func NewInternalProductHandler(cfg Config) *internalProductHandler {
	return &internalProductHandler{core: cfg.Core}
}

func (h *internalProductHandler) Register(router fiber.Router) {
	router.Get("/", h.getInternalProducts)
	router.Get("/:id", h.getInternalProduct)
	router.Post("/", h.createInternalProduct)
	router.Put("/:id", h.updateInternalProduct)
	router.Delete("/:id", h.deleteInternalProduct)

	router.Get("/:id/pricings", h.getInternalProductPricings)
	router.Post("/:id/pricings", h.createInternalProductPricing)
	router.Get("/pricings/:pricing_id", h.getInternalProductPricing)
	router.Put("/pricings/:pricing_id", h.updateInternalProductPricing)
	router.Delete("/pricings/:pricing_id", h.deleteInternalProductPricing)

	router.Get("/pricings/:pricing_id/prices", h.getInternalProductPrices)
	router.Post("/pricings/:pricing_id/prices", h.createInternalProductPrice)
	router.Put("/prices/:price_id", h.updateInternalProductPrice)
	router.Delete("/prices/:price_id", h.deleteInternalProductPrice)
}
