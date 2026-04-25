package handler

import (
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/response"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *webhookHandler) handleDOKUWebhook(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:webhook:handle_doku_webhook:handleDOKUWebhook")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx        = c.UserContext()
		rawBody    = append([]byte(nil), c.Body()...)
		targetPath = c.Path()
		req        restentity.DOKUWebhookNotification
	)

	if len(rawBody) == 0 {
		err := errors.New("request body is required")
		log.Ctx(ctx).Warn().Err(err).Msg("DOKU webhook body is empty")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := json.Unmarshal(rawBody, &req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Bytes(common.LogKeyPayload, rawBody).Msg("Invalid DOKU webhook payload")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	notification := mapper.DOKUWebhookNotificationToCore(
		targetPath,
		rawBody,
		coreentity.DOKUWebhookSignatureHeaders{
			ClientID:         c.Get("Client-Id"),
			RequestID:        c.Get("Request-Id"),
			RequestTimestamp: c.Get("Request-Timestamp"),
			Signature:        c.Get("Signature"),
		},
		req,
	)

	result, err := h.core.HandleDOKUWebhook(ctx, notification)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "invalid doku webhook signature" {
			status = fiber.StatusUnauthorized
			log.Ctx(ctx).Warn().Err(err).Str("path", targetPath).Msg("Rejected DOKU webhook")
		} else {
			log.Ctx(ctx).Error().Err(err).Str("path", targetPath).Msg("Failed to process DOKU webhook")
		}
		return c.Status(status).JSON(response.Error(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(result, "Webhook processed"))
}
