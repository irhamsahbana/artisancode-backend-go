package handler

import (
	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

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
		return c.Status(fiber.StatusForbidden).JSON(response.Error(errmsg.MessageYouDontHaveAccessToThisCompany))
	}

	tree, err := h.core.GetOrgUnitTree(ctx, userCtx.TenantID, companyID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("company_id", companyID).Msg("Failed to get org unit tree")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(tree, ""))
}
