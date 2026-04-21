package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *userInvitationHandler) createInvitation(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:userinvitation:create_invitation:createInvitation")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.CreateUserInvitationReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate request body")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	created, err := h.core.CreateInvitation(ctx, coreentity.UserInvitation{
		UserCtx:    common.GetUserContext(ctx),
		Email:      req.Email,
		RoleCode:   req.RoleCode,
		EmployeeID: req.EmployeeID,
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, req).Msg("Failed to create invitation")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(restentity.CreateUserInvitationResp{
		ID:          created.ID,
		AcceptToken: created.AcceptToken,
		ExpiresAt:   created.ExpiresAt,
	}, ""))
}
