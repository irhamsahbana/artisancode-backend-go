package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
	integrationMocks "codebase-app/internal/ports/integration/mocks"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type billingRepoStub struct{}
type currencyRepoStub struct{}

func (billingRepoStub) EnsureBillingAccount(
	ctx context.Context,
	tenantID string,
) (*coreentity.InternalBillingAccount, error) {
	return &coreentity.InternalBillingAccount{
		ID:       "account-1",
		TenantID: tenantID,
		Status:   coreentity.InternalTenantBillingAccountStatusOpen,
	}, nil
}

func (billingRepoStub) GetInvoices(
	ctx context.Context,
	filter coreentity.TenantBillingInvoiceListFilter,
) ([]coreentity.TenantBillingInvoice, error) {
	return nil, nil
}

func (billingRepoStub) CreateInvoice(
	ctx context.Context,
	data coreentity.InternalTenantInvoice,
) (*coreentity.InternalTenantInvoice, error) {
	data.ID = "invoice-1"
	data.InvoiceNumber = "TINV-001"
	return &data, nil
}

func (billingRepoStub) CreatePaymentAttempt(
	ctx context.Context,
	data coreentity.InternalTenantPaymentAttempt,
) (*coreentity.InternalTenantPaymentAttempt, error) {
	data.ID = "attempt-1"
	return &data, nil
}

func (billingRepoStub) UpdatePaymentAttemptGateway(
	ctx context.Context,
	data coreentity.InternalTenantPaymentAttempt,
) (*coreentity.InternalTenantPaymentAttempt, error) {
	return &data, nil
}

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

func (currencyRepoStub) GetProviderCurrencies(
	ctx context.Context,
	filter coreentity.InternalPaymentProviderCurrencyListFilter,
) ([]coreentity.InternalPaymentProviderCurrency, int, error) {
	return nil, 0, nil
}

func (currencyRepoStub) UpsertProviderCurrency(
	ctx context.Context,
	data coreentity.InternalPaymentProviderCurrency,
) error {
	return nil
}

func (currencyRepoStub) DeleteProviderCurrency(
	ctx context.Context,
	filter coreentity.InternalPaymentProviderCurrencyFilter,
) error {
	return nil
}

func (currencyRepoStub) IsProviderCurrencyActive(
	ctx context.Context,
	provider string,
	currencyCode string,
	amount string,
) (bool, error) {
	return true, nil
}

func TestCreateCheckoutIncludesDOKULineItemName(t *testing.T) {
	ctx := context.Background()
	productRepo := dbMocks.NewInternalProductRepository(t)
	doku := integrationMocks.NewDokuClient(t)

	productRepo.EXPECT().
		GetInternalProductPrice(
			mock.Anything,
			coreentity.InternalProductPriceFilter{ID: "price-1"},
		).
		Return(&coreentity.InternalProductPrice{
			ID:                       "price-1",
			InternalProductPricingID: "pricing-1",
			CurrencyCode:             "IDR",
			Amount:                   decimal.NewFromInt(999000),
			StartedAt:                "2026-01-01T00:00:00Z",
		}, nil)
	productRepo.EXPECT().
		GetInternalProductPricing(
			mock.Anything,
			coreentity.InternalProductPricingFilter{ID: "pricing-1"},
		).
		Return(&coreentity.InternalProductPricing{
			ID:                "pricing-1",
			InternalProductID: "product-1",
			Name:              "Annual",
			Status:            coreentity.InternalProductStatusActive,
		}, nil)
	productRepo.EXPECT().
		GetInternalProduct(
			mock.Anything,
			coreentity.InternalProductFilter{ID: "product-1"},
		).
		Return(&coreentity.InternalProduct{
			ID:     "product-1",
			Name:   "Attendance Pro",
			Status: coreentity.InternalProductStatusActive,
		}, nil)
	doku.EXPECT().
		CreatePayment(mock.Anything, mock.MatchedBy(func(req restentity.DokuCreatePaymentRequest) bool {
			return req.InvoiceNumber == "TINV-001" &&
				req.Amount == 999000 &&
				req.Currency == "IDR" &&
				req.CustomerName == "irhamsahbana" &&
				req.CustomerEmail == "billing+tenant-1@example.com" &&
				len(req.LineItems) == 1 &&
				req.LineItems[0].Name == "Attendance Pro - Annual" &&
				req.LineItems[0].Price == 999000 &&
				req.LineItems[0].Quantity == 1
		})).
		Return(&restentity.DokuCreatePaymentResponse{
			RequestID:  "req-1",
			PaymentURL: "https://pay.example.test",
		}, nil)

	core := NewInternalTenantBillingCore(Config{
		ProductRepo:  productRepo,
		CurrencyRepo: currencyRepoStub{},
		BillingRepo:  billingRepoStub{},
		DOKU:         doku,
	})

	result, err := core.CreateCheckout(ctx, coreentity.TenantBillingCheckoutInput{
		UserCtx: common.UserContext{
			UserName: "irhamsahbana",
			TenantID: "tenant-1",
			Roles:    []string{"owner"},
		},
		PriceID: "price-1",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "https://pay.example.test", result.PaymentURL)
}

func TestBuildCheckoutLineItemName(t *testing.T) {
	t.Run("combines product and pricing when pricing is generic", func(t *testing.T) {
		require.Equal(t, "Attendance Pro - Annual", buildCheckoutLineItemName("Attendance Pro", "Annual"))
	})

	t.Run("keeps pricing name when it already contains product", func(t *testing.T) {
		require.Equal(
			t,
			"Attendance Pro Annual",
			buildCheckoutLineItemName("Attendance Pro", "Attendance Pro Annual"),
		)
	})
}
