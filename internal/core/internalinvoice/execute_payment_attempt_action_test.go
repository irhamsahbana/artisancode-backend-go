package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExecutePaymentAttemptAction(t *testing.T) {
	ctx := context.Background()
	input := coreentity.InternalPaymentAttemptActionInput{
		UserCtx: common.UserContext{TenantID: "tenant-1"},
		ID:      "attempt-1",
		Action:  coreentity.ActionRetryPayment,
	}

	tests := []struct {
		name      string
		input     coreentity.InternalPaymentAttemptActionInput
		setup     func(invoiceRepo *dbMocks.InternalInvoiceRepository)
		wantError bool
	}{
		{
			name:  "retries failed attempt using previous manual provider",
			input: input,
			setup: func(invoiceRepo *dbMocks.InternalInvoiceRepository) {
				invoiceRepo.EXPECT().
					GetPaymentAttempt(mock.Anything, "tenant-1", "attempt-1").
					Return(&coreentity.InternalPaymentAttempt{
						ID:                "attempt-1",
						InternalInvoiceID: "invoice-1",
						Provider:          coreentity.PaymentProviderManual,
						Status:            coreentity.PaymentAttemptStatusFailed,
					}, nil)
				invoiceRepo.EXPECT().
					GetInvoice(mock.Anything, "tenant-1", "invoice-1").
					Return(&coreentity.InternalInvoice{
						ID:                "invoice-1",
						InvoiceNumber:     "INV-001",
						Status:            coreentity.InvoiceStatusOpen,
						AmountOutstanding: decimal.NewFromInt(100),
					}, nil)
				invoiceRepo.EXPECT().
					CreatePaymentAttempt(mock.Anything, mock.MatchedBy(func(data coreentity.InternalPaymentAttempt) bool {
						return data.InternalInvoiceID == "invoice-1" &&
							data.Provider == coreentity.PaymentProviderManual
					})).
					Return(&coreentity.InternalPaymentAttempt{ID: "attempt-2"}, nil)
			},
		},
		{
			name: "rejects unsupported action",
			input: coreentity.InternalPaymentAttemptActionInput{
				UserCtx: common.UserContext{TenantID: "tenant-1"},
				ID:      "attempt-1",
				Action:  "unsupported",
			},
			wantError: true,
		},
		{
			name:  "rejects non-retryable attempt status",
			input: input,
			setup: func(invoiceRepo *dbMocks.InternalInvoiceRepository) {
				invoiceRepo.EXPECT().
					GetPaymentAttempt(mock.Anything, "tenant-1", "attempt-1").
					Return(&coreentity.InternalPaymentAttempt{
						ID:     "attempt-1",
						Status: coreentity.PaymentAttemptStatusPending,
					}, nil)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoiceRepo := dbMocks.NewInternalInvoiceRepository(t)
			if tt.setup != nil {
				tt.setup(invoiceRepo)
			}
			core := NewInternalInvoiceCore(Config{InvoiceRepo: invoiceRepo})

			got, err := core.ExecutePaymentAttemptAction(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, "attempt-2", got.PaymentAttempt.ID)
		})
	}
}
