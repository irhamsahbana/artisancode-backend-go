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

func (h *userInvitationHandler) resendInvitation(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:userinvitation:resend_invitation:resendInvitation")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.ResendUserInvitationReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse route params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate route params")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	item, err := h.core.ResendInvitation(ctx, coreentity.UserInvitation{
		UserCtx: common.GetUserContext(ctx),
		ID:      req.ID,
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, req).Msg(errmsg.MessageFailedToResendInvitation)
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(restentity.ResendUserInvitationResp{
		ID:          item.ID,
		AcceptToken: item.AcceptToken,
		ExpiresAt:   item.ExpiresAt,
		EmailSent:   item.EmailSent,
	}, invitationResponseMessage(item.EmailSent)))
}
