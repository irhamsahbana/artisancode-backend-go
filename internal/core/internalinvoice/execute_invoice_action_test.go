package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
	infraConfig "codebase-app/internal/infrastructure/config"
	integrationMocks "codebase-app/internal/ports/integration/mocks"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	infraConfig.Envs = &infraConfig.Config{}
	infraConfig.Envs.App.Name = "internalinvoice-core-test"
}

type currencyRepoStub struct{}

func (currencyRepoStub) GetInternalCurrencies(
	ctx context.Context,
	filter coreentity.InternalCurrencyListFilter,
) ([]coreentity.InternalCurrency, int, error) {
	return nil, 0, nil
}

func (currencyRepoStub) GetInternalCurrency(
	ctx context.Context,
	filter coreentity.InternalCurrencyFilter,
) (*coreentity.InternalCurrency, error) {
	return &coreentity.InternalCurrency{
		Code:          filter.Code,
		Symbol:        "Rp",
		DecimalPlaces: 0,
		IsActive:      true,
		IsDefault:     true,
	}, nil
}

func (currencyRepoStub) CreateInternalCurrency(
	ctx context.Context,
	data coreentity.InternalCurrency,
) (*coreentity.InternalCurrency, error) {
	return nil, nil
}

func (currencyRepoStub) UpdateInternalCurrency(
	ctx context.Context,
	data coreentity.InternalCurrency,
) error {
	return nil
}

func (currencyRepoStub) DeleteInternalCurrency(
	ctx context.Context,
	filter coreentity.InternalCurrencyDeleteFilter,
) error {
	return nil
}

func (currencyRepoStub) IsCurrencyActive(ctx context.Context, code string) (bool, error) {
	return true, nil
}

func (currencyRepoStub) GetDefaultCurrency(ctx context.Context) (*coreentity.InternalCurrency, error) {
	return &coreentity.InternalCurrency{Code: "IDR", IsActive: true, IsDefault: true}, nil
}

func TestExecuteInvoiceAction(t *testing.T) {
	ctx := context.Background()
	input := coreentity.InternalInvoiceActionInput{
		UserCtx: common.UserContext{TenantID: "tenant-1"},
		ID:      "invoice-1",
		Action:  coreentity.ActionCreateCustomPaymentAttempt,
	}

	tests := []struct {
		name      string
		input     coreentity.InternalInvoiceActionInput
		doku      bool
		setup     func(invoiceRepo *dbMocks.InternalInvoiceRepository, doku *integrationMocks.DokuClient)
		wantError bool
	}{
		{
			name:  "creates manual payment attempt with defaults",
			input: input,
			setup: func(invoiceRepo *dbMocks.InternalInvoiceRepository, doku *integrationMocks.DokuClient) {
				invoiceRepo.EXPECT().
					GetInvoice(mock.Anything, "tenant-1", "invoice-1").
					Return(&coreentity.InternalInvoice{
						ID:                "invoice-1",
						InvoiceNumber:     "INV-001",
						Status:            coreentity.InvoiceStatusOpen,
						AmountOutstanding: decimal.NewFromInt(250000),
					}, nil)
				invoiceRepo.EXPECT().
					CreatePaymentAttempt(mock.Anything, mock.MatchedBy(func(data coreentity.InternalPaymentAttempt) bool {
						return data.InternalInvoiceID == "invoice-1" &&
							data.Provider == coreentity.PaymentProviderManual &&
							data.PaymentMethodType == "bank_transfer" &&
							data.PaymentChannelCode == "manual_bank_transfer" &&
							data.Status == coreentity.PaymentAttemptStatusPending &&
							data.RequestedAmount.Equal(decimal.NewFromInt(250000))
					})).
					Return(&coreentity.InternalPaymentAttempt{ID: "attempt-1", Provider: coreentity.PaymentProviderManual}, nil)
			},
		},
		{
			name: "creates DOKU gateway attempt",
			input: coreentity.InternalInvoiceActionInput{
				UserCtx:       common.UserContext{TenantID: "tenant-1"},
				ID:            "invoice-1",
				Action:        coreentity.ActionCreateDOKUPaymentAttempt,
				CustomerName:  "Budi",
				CustomerEmail: "budi@example.com",
			},
			doku: true,
			setup: func(invoiceRepo *dbMocks.InternalInvoiceRepository, doku *integrationMocks.DokuClient) {
				invoiceRepo.EXPECT().
					GetInvoice(mock.Anything, "tenant-1", "invoice-1").
					Return(&coreentity.InternalInvoice{
						ID:                "invoice-1",
						InvoiceNumber:     "INV-001",
						Status:            coreentity.InvoiceStatusOpen,
						CurrencyCode:      "IDR",
						AmountOutstanding: decimal.NewFromInt(250000),
					}, nil)
				invoiceRepo.EXPECT().
					CreatePaymentAttempt(mock.Anything, mock.MatchedBy(func(data coreentity.InternalPaymentAttempt) bool {
						return data.Provider == coreentity.PaymentProviderDOKU &&
							data.Status == coreentity.PaymentAttemptStatusInitiated
					})).
					Return(&coreentity.InternalPaymentAttempt{ID: "attempt-1", InternalInvoiceID: "invoice-1"}, nil)
				doku.EXPECT().
					CreatePayment(mock.Anything, mock.MatchedBy(func(req restentity.DokuCreatePaymentRequest) bool {
						return req.InvoiceNumber == "INV-001" &&
							req.Amount == 250000 &&
							req.CustomerName == "Budi" &&
							req.CustomerEmail == "budi@example.com"
					})).
					Return(&restentity.DokuCreatePaymentResponse{RequestID: "req-1", PaymentURL: "https://pay.test"}, nil)
				invoiceRepo.EXPECT().
					UpdatePaymentAttemptGateway(mock.Anything, mock.MatchedBy(func(data coreentity.InternalPaymentAttempt) bool {
						return data.ProviderRequestID == "req-1" &&
							data.ProviderPaymentURL == "https://pay.test" &&
							data.Status == coreentity.PaymentAttemptStatusPending
					})).
					Return(&coreentity.InternalPaymentAttempt{ID: "attempt-1", Status: coreentity.PaymentAttemptStatusPending}, nil)
			},
		},
		{
			name: "rejects unsupported action",
			input: coreentity.InternalInvoiceActionInput{
				UserCtx: common.UserContext{TenantID: "tenant-1"},
				ID:      "invoice-1",
				Action:  "unsupported",
			},
			wantError: true,
		},
		{
			name: "rejects DOKU action when client missing",
			input: coreentity.InternalInvoiceActionInput{
				UserCtx: common.UserContext{TenantID: "tenant-1"},
				ID:      "invoice-1",
				Action:  coreentity.ActionCreateDOKUPaymentAttempt,
			},
			wantError: true,
		},
		{
			name:  "rejects non-payable invoice",
			input: input,
			setup: func(invoiceRepo *dbMocks.InternalInvoiceRepository, doku *integrationMocks.DokuClient) {
				invoiceRepo.EXPECT().
					GetInvoice(mock.Anything, "tenant-1", "invoice-1").
					Return(&coreentity.InternalInvoice{ID: "invoice-1", Status: coreentity.InvoiceStatusPaid}, nil)
			},
			wantError: true,
		},
		{
			name:  "returns repository error",
			input: input,
			setup: func(invoiceRepo *dbMocks.InternalInvoiceRepository, doku *integrationMocks.DokuClient) {
				invoiceRepo.EXPECT().
					GetInvoice(mock.Anything, "tenant-1", "invoice-1").
					Return(nil, errors.New("repo failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoiceRepo := dbMocks.NewInternalInvoiceRepository(t)
			var doku *integrationMocks.DokuClient
			if tt.doku {
				doku = integrationMocks.NewDokuClient(t)
			}
			if tt.setup != nil {
				tt.setup(invoiceRepo, doku)
			}
			cfg := Config{InvoiceRepo: invoiceRepo}
			if tt.doku {
				cfg.DOKU = doku
				cfg.CurrencyRepo = currencyRepoStub{}
			}
			core := NewInternalInvoiceCore(cfg)

			got, err := core.ExecuteInvoiceAction(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, got.PaymentAttempt.ID)
		})
	}
}
