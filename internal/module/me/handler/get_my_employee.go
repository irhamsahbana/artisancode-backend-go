package handler

import (
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *meHandler) getMyEmployee(c *fiber.Ctx) error {
	ctx := c.UserContext()
	uc := common.GetUserContext(ctx)

	item, err := h.core.GetMyEmployee(ctx, coreentity.SelfFilter{
		UserCtx:  uc,
		TenantID: uc.TenantID,
		UserID:   uc.UserID,
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get current employee profile")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(mapper.MyEmployeeFromCoreToRest(*item), ""))
}
