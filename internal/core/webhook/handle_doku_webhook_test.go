package core

import (
	"codebase-app/internal/entity/coreentity"
	infraConfig "codebase-app/internal/infrastructure/config"
	integrationMocks "codebase-app/internal/ports/integration/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	infraConfig.Envs = &infraConfig.Config{}
	infraConfig.Envs.App.Name = "webhook-core-test"
}

func TestHandleDOKUWebhook(t *testing.T) {
	ctx := context.Background()
	notification := coreentity.DOKUWebhookNotification{
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
	}
	verifier := integrationMocks.NewDOKUWebhookVerifier(t)
	verifier.EXPECT().
		VerifyWebhookSignatureHeaders(notification.Headers, notification.RawBody, notification.TargetPath).
		Return(true)

	core := NewWebhookCore(Config{
		DOKUVerifier: verifier,
	})

	result, err := core.HandleDOKUWebhook(ctx, notification)

	require.NoError(t, err)
	require.Equal(t, "doku", result.Provider)
	require.Equal(t, "INV-001", result.InvoiceNumber)
	require.Equal(t, "PAID", result.OrderStatus)
	require.Equal(t, "SUCCESS", result.PaymentStatus)
}

func TestHandleDOKUWebhookInvalidSignature(t *testing.T) {
	notification := coreentity.DOKUWebhookNotification{
		TargetPath: "/webhooks/doku",
		RawBody:    []byte(`{}`),
	}
	verifier := integrationMocks.NewDOKUWebhookVerifier(t)
	verifier.EXPECT().
		VerifyWebhookSignatureHeaders(notification.Headers, notification.RawBody, notification.TargetPath).
		Return(false)

	core := NewWebhookCore(Config{
		DOKUVerifier: verifier,
	})

	_, err := core.HandleDOKUWebhook(context.Background(), notification)

	require.Error(t, err)
}

func TestHandleDOKUWebhookMissingVerifier(t *testing.T) {
	core := NewWebhookCore(Config{})

	_, err := core.HandleDOKUWebhook(context.Background(), coreentity.DOKUWebhookNotification{})

	require.Error(t, err)
}
