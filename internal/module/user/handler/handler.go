package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	corePorts "codebase-app/internal/ports/core"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type userHandler struct {
	core corePorts.UserCore
}

type UserHandlerConfig struct {
	Core corePorts.UserCore
}

func NewUserHandler(cfg UserHandlerConfig) *userHandler {
	return &userHandler{
		core: cfg.Core,
	}
}

func (h *userHandler) Register(router fiber.Router) {
	router.Post("/login", h.login)
	router.Post("/register", h.register)
	router.Post("/refresh-token", h.refreshToken)
}

func (h *userHandler) login(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.LoginReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request validation")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	user := mapper.LoginReqToCore(ctx, *req)
	tokens, err := h.core.Login(ctx, user)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"email": req.Email}).Msg("Login service error")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp := mapper.AuthTokensToLoginResp(*tokens)
	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *userHandler) register(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.RegisterReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request validation")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	user := mapper.RegisterReqToCore(ctx, *req)
	tenant := mapper.RegisterReqToTenant(ctx, *req)
	tokens, err := h.core.RegisterOwner(ctx, user, tenant)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, req).Msg("Register service error")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp := mapper.AuthTokensToRegisterResp(*tokens)
	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *userHandler) refreshToken(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.RefreshTokenReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request validation")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	user := mapper.RefreshTokenReqToCore(ctx, *req)
	tokens, err := h.core.RefreshToken(ctx, user)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("RefreshToken service error")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp := mapper.AuthTokensToLoginResp(*tokens)
	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
