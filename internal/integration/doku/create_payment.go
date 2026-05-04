package doku

import (
	"bytes"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

func (c *dokuClient) CreatePayment(ctx context.Context, req restentity.DokuCreatePaymentRequest) (*restentity.DokuCreatePaymentResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:integration:doku:create_payment:CreatePayment")
	defer span.End()

	if err := c.validate(); err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}
	if err := validateCreatePaymentRequest(req); err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	callbackURL := strings.TrimSpace(req.CallbackURL)
	if callbackURL == "" {
		callbackURL = c.callbackURL
	}
	if callbackURL == "" {
		err := errors.New("doku callback url is required")
		tracing.RecordError(span, err)
		return nil, err
	}

	targetPath := "/checkout/v1/payment"
	requestID := c.newRequestID()
	timestamp := c.now().UTC().Format("2006-01-02T15:04:05Z")
	payload := restentity.DokuCheckoutRequest{
		Order: restentity.DokuCheckoutOrder{
			Amount:        req.Amount,
			InvoiceNumber: req.InvoiceNumber,
			Currency:      req.Currency,
			CallbackURL:   callbackURL,
			AutoRedirect:  req.AutoRedirectOrDefault(),
			LineItems:     req.LineItems,
		},
		Payment: restentity.DokuCheckoutPayment{
			PaymentDueDate: req.ExpiryMinutesOrDefault(),
		},
		Customer: restentity.DokuCheckoutCustomer{
			Name:    req.CustomerName,
			Email:   req.CustomerEmail,
			Phone:   req.CustomerPhone,
			Address: req.CustomerAddress,
			Country: req.CustomerCountryOrDefault(),
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+targetPath, bytes.NewReader(payloadBytes))
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Client-Id", c.clientID)
	httpReq.Header.Set("Request-Id", requestID)
	httpReq.Header.Set("Request-Timestamp", timestamp)
	httpReq.Header.Set("Signature", c.generateSignature(payloadBytes, timestamp, requestID, targetPath))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("invoice_number", req.InvoiceNumber).Msg("failed to call DOKU checkout API")
		tracing.RecordError(span, err)
		return nil, err
	}
	defer resp.Body.Close()

	var body restentity.DokuCheckoutResponseEnvelope
	if err := decodeResponse(resp, &body); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("invoice_number", req.InvoiceNumber).Msg("failed to decode DOKU checkout API response")
		tracing.RecordError(span, err)
		return nil, err
	}

	return &restentity.DokuCreatePaymentResponse{
		InvoiceID:     body.Response.Order.InvoiceNumber,
		InvoiceNumber: body.Response.Order.InvoiceNumber,
		Amount:        int64(body.Response.Order.Amount),
		PaymentURL:    body.Response.Payment.URL,
		RequestID:     requestID,
		Message:       body.Message,
	}, nil
}

func validateCreatePaymentRequest(req restentity.DokuCreatePaymentRequest) error {
	switch {
	case strings.TrimSpace(req.InvoiceNumber) == "":
		return errors.New("invoice number is required")
	case req.Amount <= 0:
		return errors.New("amount must be greater than zero")
	case len(strings.TrimSpace(req.Currency)) != 3 || strings.ToUpper(strings.TrimSpace(req.Currency)) != strings.TrimSpace(req.Currency):
		return errors.New("currency must be a 3-letter uppercase code")
	case strings.TrimSpace(req.CustomerEmail) == "":
		return errors.New("customer email is required")
	case strings.TrimSpace(req.CustomerName) == "":
		return errors.New("customer name is required")
	default:
		return nil
	}
}
