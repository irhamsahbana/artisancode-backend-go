package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	portsCore "codebase-app/internal/ports/core"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"codebase-app/pkg/types"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type employeeHandler struct {
	core portsCore.EmployeeCore
}

func NewEmployeeHandler(core portsCore.EmployeeCore) *employeeHandler {
	return &employeeHandler{core: core}
}

func (h *employeeHandler) Register(router fiber.Router) {
	router.Get("/", h.getEmployees)
	router.Get("/:id", h.getEmployee)
	router.Post("/", h.createEmployee)
	router.Put("/:id", h.updateEmployee)
	router.Delete("/:id", h.deleteEmployee)
}

func (h *employeeHandler) getEmployees(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.GetEmployeesReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to validate query params")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	filter := coreentity.EmployeeListFilter{
		TenantID:  common.GetUserContext(ctx).TenantID,
		Q:         req.Q,
		Status:    req.Status,
		OrgUnitID: req.OrgUnitID,
		Page:      req.Page,
		Paginate:  req.Paginate,
	}

	items, total, err := h.core.GetEmployees(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get employees")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restItems := make([]restentity.Employee, 0, len(items))
	for _, item := range items {
		restItems = append(restItems, mapper.EmployeeFromCoreToRest(item))
	}

	resp := &restentity.GetEmployeesResp{
		Items: restItems,
		Meta: types.Meta{
			Page:      req.Page,
			Paginate:  req.Paginate,
			TotalData: total,
		},
	}
	resp.Meta.CountTotalPage(req.Page, req.Paginate, total)

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *employeeHandler) getEmployee(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.GetEmployeeReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := c.QueryParser(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to validate params")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	filter := coreentity.Employee{
		TenantID: common.GetUserContext(ctx).TenantID,
		ID:       req.ID,
	}

	item, err := h.core.GetEmployee(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get employee")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(restentity.GetEmployeeResp{
		Employee: mapper.EmployeeFromCoreToRest(*item),
	}, ""))
}

func (h *employeeHandler) createEmployee(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.CreateEmployeeReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to validate request body")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	data := mapper.EmployeeFromRestCreateToCore(ctx, *req)
	created, err := h.core.CreateEmployee(ctx, data)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create employee")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(restentity.CreateEmployeeResp{ID: created.ID}, ""))
}

func (h *employeeHandler) updateEmployee(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.UpdateEmployeeReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := c.QueryParser(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to validate request body")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	data := mapper.EmployeeFromRestUpdateToCore(ctx, *req)
	if err := h.core.UpdateEmployee(ctx, data); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to update employee")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}

func (h *employeeHandler) deleteEmployee(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.DeleteEmployeeReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := c.QueryParser(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to validate params")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	filter := coreentity.EmployeeDeleteFilter{
		TenantID: common.GetUserContext(ctx).TenantID,
		ID:       req.ID,
	}

	if err := h.core.DeleteEmployee(ctx, filter); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete employee")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}