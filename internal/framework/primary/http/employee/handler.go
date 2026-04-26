package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type employeeHandler struct {
	core portsCore.EmployeeCore
}

type Config struct {
	Core portsCore.EmployeeCore
}

func NewEmployeeHandler(cfg Config) *employeeHandler {
	return &employeeHandler{core: cfg.Core}
}

func (h *employeeHandler) Register(router fiber.Router) {
	router.Get("/", h.getEmployees)
	router.Get("/:id", h.getEmployee)
	router.Post("/", h.createEmployee)
	router.Put("/:id", h.updateEmployee)
	router.Delete("/:id", h.deleteEmployee)
}
