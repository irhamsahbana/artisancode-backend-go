package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"codebase-app/pkg/types"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *userHandler) getUsers(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:user:get_users:getUsers")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.GetUsersReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate query params")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	uc := common.GetUserContext(ctx)
	filter := coreentity.UserListFilter{
		TenantID: uc.TenantID,
		Q:        req.Q,
		Page:     req.Page,
		Paginate: req.Limit,
	}

	items, total, err := h.core.GetUsers(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get users")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restItems := make([]restentity.UserResource, 0, len(items))
	for _, item := range items {
		restItems = append(restItems, mapper.UserFromCoreToRest(item))
	}

	meta := types.Meta{}
	meta.CountTotalPage(req.Page, req.Limit, total)

	resp := &restentity.GetUsersResp{
		Items:      restItems,
		Pagination: restentity.NewPaginationResp(meta),
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
