package integration

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
	"context"
	"net/http"
)

type DOKUWebhookVerifier interface {
	VerifyWebhookSignatureHeaders(headers coreentity.DOKUWebhookSignatureHeaders, body []byte, targetPath string) bool
}

type DokuClient interface {
	CreatePayment(ctx context.Context, req restentity.DokuCreatePaymentRequest) (*restentity.DokuCreatePaymentResponse, error)
	CheckStatus(ctx context.Context, invoiceNumber string) (*restentity.DokuCheckStatusResponse, error)
	VerifyWebhookSignature(headers http.Header, body []byte, targetPath string) bool
	VerifyWebhookSignatureHeaders(headers coreentity.DOKUWebhookSignatureHeaders, body []byte, targetPath string) bool
	DecodedPublicKey() string
}
