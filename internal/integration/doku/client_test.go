package doku

import (
	"codebase-app/internal/entity/restentity"
	infraConfig "codebase-app/internal/infrastructure/config"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func init() {
	infraConfig.Envs = &infraConfig.Config{}
	infraConfig.Envs.App.Name = "doku-test"
}

func TestCreatePayment(t *testing.T) {
	const (
		expectedRequestID = "req-123"
		expectedTimestamp = "2026-04-24T10:11:12Z"
	)

	httpClient := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST method, got %s", r.Method)
			}
			if r.URL.Path != "/checkout/v1/payment" {
				t.Fatalf("expected checkout path, got %s", r.URL.Path)
			}
			if got := r.Header.Get("Client-Id"); got != "client-123" {
				t.Fatalf("expected client id header, got %s", got)
			}
			if got := r.Header.Get("Request-Id"); got != expectedRequestID {
				t.Fatalf("expected request id %s, got %s", expectedRequestID, got)
			}
			if got := r.Header.Get("Request-Timestamp"); got != expectedTimestamp {
				t.Fatalf("expected request timestamp %s, got %s", expectedTimestamp, got)
			}
			if got := r.Header.Get("Signature"); got == "" {
				t.Fatal("expected signature header")
			}

			var payload restentity.DokuCheckoutRequest
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode payload: %v", err)
			}
			if payload.Order.CallbackURL != "https://example.com/webhooks/doku" {
				t.Fatalf("unexpected callback url: %s", payload.Order.CallbackURL)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"message":["SUCCESS"],"response":{"order":{"invoice_number":"INV-001","amount":15000},"payment":{"url":"https://pay.example.com/abc"}}}`)),
			}, nil
		}),
	}

	dokuClient := NewClient(Config{
		BaseURL:     "https://api-sandbox.doku.com",
		ClientID:    "client-123",
		SecretKey:   "secret-123",
		CallbackURL: "https://example.com/webhooks/doku",
		HTTPClient:  httpClient,
		Now: func() time.Time {
			return time.Date(2026, 4, 24, 10, 11, 12, 0, time.UTC)
		},
		NewRequestID: func() string {
			return expectedRequestID
		},
	})

	resp, err := dokuClient.CreatePayment(context.Background(), restentity.DokuCreatePaymentRequest{
		InvoiceNumber: "INV-001",
		Amount:        15000,
		Currency:      "IDR",
		CustomerEmail: "buyer@example.com",
		CustomerName:  "Buyer Test",
	})
	if err != nil {
		t.Fatalf("CreatePayment returned error: %v", err)
	}
	if resp.PaymentURL != "https://pay.example.com/abc" {
		t.Fatalf("expected payment url, got %s", resp.PaymentURL)
	}
	if resp.RequestID != expectedRequestID {
		t.Fatalf("expected request id %s, got %s", expectedRequestID, resp.RequestID)
	}
}

func TestCheckStatus(t *testing.T) {
	httpClient := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET method, got %s", r.Method)
			}
			if r.URL.Path != "/orders/v1/status/INV-001" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"order":{"invoice_number":"INV-001","amount":15000,"status":"PAID"},"transaction":{"status":"SUCCESS","date":"2026-04-24T10:20:00Z","original_request_id":"req-123"}}`)),
			}, nil
		}),
	}

	dokuClient := NewClient(Config{
		BaseURL:    "https://api-sandbox.doku.com",
		ClientID:   "client-123",
		SecretKey:  "secret-123",
		HTTPClient: httpClient,
	})

	resp, err := dokuClient.CheckStatus(context.Background(), "INV-001")
	if err != nil {
		t.Fatalf("CheckStatus returned error: %v", err)
	}
	if resp.Order.Status != "PAID" {
		t.Fatalf("expected PAID status, got %s", resp.Order.Status)
	}
}

func TestVerifyWebhookSignature(t *testing.T) {
	publicKey := base64.StdEncoding.EncodeToString([]byte("PUBLIC-KEY"))
	dokuClient := NewClient(Config{
		BaseURL:    "https://api-sandbox.doku.com",
		ClientID:   "client-123",
		SecretKey:  "secret-123",
		PublicKey:  publicKey,
		HTTPClient: &http.Client{},
	})

	body := []byte(`{"order":{"invoice_number":"INV-001"}}`)
	targetPath := "/webhooks/doku"
	timestamp := "2026-04-24T10:11:12Z"
	requestID := "req-123"

	signature := dokuClient.generateSignature(body, timestamp, requestID, targetPath)

	headers := http.Header{}
	headers.Set("Client-Id", "client-123")
	headers.Set("Request-Id", requestID)
	headers.Set("Request-Timestamp", timestamp)
	headers.Set("Signature", signature)

	if !dokuClient.VerifyWebhookSignature(headers, body, targetPath) {
		t.Fatal("expected webhook signature to be valid")
	}
	if dokuClient.DecodedPublicKey() != "PUBLIC-KEY" {
		t.Fatalf("expected decoded public key, got %s", dokuClient.DecodedPublicKey())
	}
}
