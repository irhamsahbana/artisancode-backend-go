package doku

import (
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func (c *dokuClient) CancelOrder(ctx context.Context, invoiceNumber string) (*restentity.DokuCancelOrderResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:integration:doku:cancel_order:CancelOrder")
	defer span.End()

	if err := c.validate(); err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}
	if strings.TrimSpace(invoiceNumber) == "" {
		err := errors.New("invoice number is required")
		tracing.RecordError(span, err)
		return nil, err
	}

	targetPath := "/checkout/v1/order/" + invoiceNumber
	requestID := c.newRequestID()
	timestamp := c.now().UTC().Format("2006-01-02T15:04:05Z")

	var nilBody []byte
	_ = nilBody
	payload, err := json.Marshal(map[string]string{"invoice_number": invoiceNumber})
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+targetPath, strings.NewReader(string(payload)))
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Client-Id", c.clientID)
	httpReq.Header.Set("Request-Id", requestID)
	httpReq.Header.Set("Request-Timestamp", timestamp)
	httpReq.Header.Set("Signature", c.generateSignature(payload, timestamp, requestID, targetPath))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}
	defer resp.Body.Close()

	var body restentity.DokuCancelOrderResponse
	if err := decodeResponse(resp, &body); err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	return &body, nil
}
