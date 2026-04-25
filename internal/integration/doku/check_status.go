package doku

import (
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

func (c *dokuClient) CheckStatus(ctx context.Context, invoiceNumber string) (*restentity.DokuCheckStatusResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:integration:doku:check_status:CheckStatus")
	defer span.End()

	if err := c.validate(); err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}
	invoiceNumber = strings.TrimSpace(invoiceNumber)
	if invoiceNumber == "" {
		err := errors.New("invoice number is required")
		tracing.RecordError(span, err)
		return nil, err
	}

	targetPath := "/orders/v1/status/" + url.PathEscape(invoiceNumber)
	requestID := c.newRequestID()
	timestamp := c.now().UTC().Format("2006-01-02T15:04:05Z")

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+targetPath, nil)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	httpReq.Header.Set("Client-Id", c.clientID)
	httpReq.Header.Set("Request-Id", requestID)
	httpReq.Header.Set("Request-Timestamp", timestamp)
	httpReq.Header.Set("Signature", c.generateGetSignature(timestamp, requestID, targetPath))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("invoice_number", invoiceNumber).Msg("failed to call DOKU status API")
		tracing.RecordError(span, err)
		return nil, err
	}
	defer resp.Body.Close()

	var body restentity.DokuCheckStatusResponse
	if err := decodeResponse(resp, &body); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("invoice_number", invoiceNumber).Msg("failed to decode DOKU status API response")
		tracing.RecordError(span, err)
		return nil, err
	}

	return &body, nil
}
