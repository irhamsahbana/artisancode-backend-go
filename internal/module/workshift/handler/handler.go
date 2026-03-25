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

type workShiftHandler struct {
	core portsCore.WorkShiftCore
}

type WorkShiftHandlerConfig struct {
	Core portsCore.WorkShiftCore
}

func NewWorkShiftHandler(cfg WorkShiftHandlerConfig) *workShiftHandler {
	return &workShiftHandler{core: cfg.Core}
}

func (h *workShiftHandler) Register(router fiber.Router) {
	router.Get("/", h.getWorkShifts)
	router.Get("/:id", h.getWorkShift)
	router.Post("/", h.createWorkShift)
	router.Put("/:id", h.updateWorkShift)
	router.Delete("/:id", h.deleteWorkShift)
}

func (h *workShiftHandler) getWorkShifts(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.GetWorkShiftsReq)
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

	filter := coreentity.WorkShiftListFilter{
		TenantID:  common.GetUserContext(ctx).TenantID,
		Q:         req.Q,
		Page:      req.Page,
		Paginate:  req.Paginate,
	}

	items, total, err := h.core.GetWorkShifts(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get work shifts")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restItems := make([]restentity.WorkShift, 0, len(items))
	for _, item := range items {
		restItems = append(restItems, mapper.WorkShiftFromCoreToRest(item))
	}

	resp := &restentity.GetWorkShiftsResp{
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

func (h *workShiftHandler) getWorkShift(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.GetWorkShiftReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to validate params")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	filter := coreentity.WorkShift{
		TenantID: common.GetUserContext(ctx).TenantID,
		ID:       req.ID,
	}

	item, err := h.core.GetWorkShift(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get work shift")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(restentity.GetWorkShiftResp{
		WorkShift: mapper.WorkShiftFromCoreToRest(*item),
	}, ""))
}

func (h *workShiftHandler) createWorkShift(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.CreateWorkShiftReq)
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

	data := mapper.WorkShiftFromRestCreateToCore(ctx, *req)
	created, err := h.core.CreateWorkShift(ctx, data)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create work shift")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(restentity.CreateWorkShiftResp{ID: created.ID}, ""))
}

func (h *workShiftHandler) updateWorkShift(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.UpdateWorkShiftReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse params")
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

	data := mapper.WorkShiftFromRestUpdateToCore(ctx, *req)
	if err := h.core.UpdateWorkShift(ctx, data); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to update work shift")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}

func (h *workShiftHandler) deleteWorkShift(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.DeleteWorkShiftReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to validate params")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	filter := coreentity.WorkShiftDeleteFilter{
		TenantID: common.GetUserContext(ctx).TenantID,
		ID:       req.ID,
	}

	if err := h.core.DeleteWorkShift(ctx, filter); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete work shift")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}