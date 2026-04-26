package core

import (
	"codebase-app/internal/entity/coreentity"
	infraConfig "codebase-app/internal/infrastructure/config"
	"context"
	"testing"
)

type fakeDOKUVerifier struct {
	valid bool
}

func (f fakeDOKUVerifier) VerifyWebhookSignatureHeaders(
	headers coreentity.DOKUWebhookSignatureHeaders,
	body []byte,
	targetPath string,
) bool {
	return f.valid
}

func init() {
	infraConfig.Envs = &infraConfig.Config{}
	infraConfig.Envs.App.Name = "webhook-core-test"
}

func TestHandleDOKUWebhook(t *testing.T) {
	core := NewWebhookCore(Config{
		DOKUVerifier: fakeDOKUVerifier{valid: true},
	})

	result, err := core.HandleDOKUWebhook(context.Background(), coreentity.DOKUWebhookNotification{
		TargetPath: "/webhooks/doku",
		RawBody:    []byte(`{"order":{"invoice_number":"INV-001"}}`),
		Headers: coreentity.DOKUWebhookSignatureHeaders{
			RequestID: "req-1",
		},
		Event: coreentity.DOKUWebhookEvent{
			Order: coreentity.DOKUWebhookOrder{
				InvoiceNumber: "INV-001",
				Status:        "PAID",
			},
			Transaction: coreentity.DOKUWebhookTransaction{
				Status: "SUCCESS",
			},
		},
	})
	if err != nil {
		t.Fatalf("HandleDOKUWebhook returned error: %v", err)
	}
	if result.Provider != "doku" {
		t.Fatalf("expected provider doku, got %s", result.Provider)
	}
}

func TestHandleDOKUWebhookInvalidSignature(t *testing.T) {
	core := NewWebhookCore(Config{
		DOKUVerifier: fakeDOKUVerifier{valid: false},
	})

	_, err := core.HandleDOKUWebhook(context.Background(), coreentity.DOKUWebhookNotification{
		TargetPath: "/webhooks/doku",
		RawBody:    []byte(`{}`),
	})
	if err == nil {
		t.Fatal("expected invalid signature error")
	}
}
