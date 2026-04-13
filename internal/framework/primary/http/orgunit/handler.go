package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	orgUnitPorts "codebase-app/internal/ports/core"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"codebase-app/pkg/types"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type orgUnitHandler struct {
	core orgUnitPorts.OrgUnitCore
}

type OrgUnitHandlerConfig struct {
	Core orgUnitPorts.OrgUnitCore
}

func NewOrgUnitHandler(cfg OrgUnitHandlerConfig) *orgUnitHandler {
	return &orgUnitHandler{core: cfg.Core}
}

func (h *orgUnitHandler) Register(router fiber.Router) {
	router.Get("/", h.getOrgUnits)
	router.Get("/tree/:companyId", h.getOrgUnitTree)
	router.Get("/:id", h.getOrgUnit)
	router.Post("/", h.createOrgUnit)
	router.Put("/:id", h.updateOrgUnit)
	router.Delete("/:id", h.deleteOrgUnit)
}

func (h *orgUnitHandler) getOrgUnitTree(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:orgunit:handler:getOrgUnitTree")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx       = c.UserContext()
		companyID = c.Params("companyId")
		userCtx   = common.GetUserContext(ctx)
	)

	if !userCtx.CanAccessCompany(companyID) {
		return c.Status(fiber.StatusForbidden).JSON(response.Error("You don't have access to this company"))
	}

	tree, err := h.core.GetOrgUnitTree(ctx, userCtx.TenantID, companyID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("company_id", companyID).Msg("Failed to get org unit tree")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(tree, ""))
}

func (h *orgUnitHandler) getOrgUnits(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:orgunit:handler:getOrgUnits")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.GetOrgUnitsReq)
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

	filter := coreentity.OrgUnitListFilter{
		TenantID: common.GetUserContext(ctx).TenantID,
		Q:        req.Q,
		Category: req.Category,
		Page:     req.Page,
		Paginate: req.Paginate,
	}

	items, total, err := h.core.GetOrgUnits(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get org units")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restItems := make([]restentity.OrgUnit, 0, len(items))
	for _, item := range items {
		restItems = append(restItems, mapper.OrgUnitFromCoreToRest(item))
	}

	resp := &restentity.GetOrgUnitsResp{
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

func (h *orgUnitHandler) getOrgUnit(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:orgunit:handler:getOrgUnit")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.GetOrgUnitReq)
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

	filter := coreentity.OrgUnit{
		TenantID: common.GetUserContext(ctx).TenantID,
		ID:       req.ID,
	}

	item, err := h.core.GetOrgUnit(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get org unit")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(restentity.GetOrgUnitResp{
		OrgUnit: mapper.OrgUnitFromCoreToRest(*item),
	}, ""))
}

func (h *orgUnitHandler) createOrgUnit(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:orgunit:handler:createOrgUnit")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.CreateOrgUnitReq)
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

	data := mapper.OrgUnitFromRestCreateToCore(ctx, *req)
	created, err := h.core.CreateOrgUnit(ctx, data)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create org unit")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(restentity.CreateOrgUnitResp{ID: created.ID}, ""))
}

func (h *orgUnitHandler) updateOrgUnit(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:orgunit:handler:updateOrgUnit")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.UpdateOrgUnitReq)
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

	data := mapper.OrgUnitFromRestUpdateToCore(ctx, *req)
	if err := h.core.UpdateOrgUnit(ctx, data); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to update org unit")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}

func (h *orgUnitHandler) deleteOrgUnit(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:orgunit:handler:deleteOrgUnit")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.DeleteOrgUnitReq)
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

	filter := coreentity.OrgUnitDeleteFilter{
		TenantID: common.GetUserContext(ctx).TenantID,
		ID:       req.ID,
	}

	if err := h.core.DeleteOrgUnit(ctx, filter); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete org unit")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}
