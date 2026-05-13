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

func (billingRepoStub) GetSubscription(
	ctx context.Context,
	tenantID string,
) (*coreentity.InternalTenantSubscription, error) {
	return &coreentity.InternalTenantSubscription{
		ID:     "sub-1",
		Status: coreentity.InternalTenantSubscriptionStatusActive,
	}, nil
}

func (billingRepoStub) GetLatestEntitlementSnapshot(
	ctx context.Context,
	tenantID string,
) (*coreentity.InternalEntitlementSnapshot, error) {
	return nil, nil
}

func (billingRepoStub) GetActiveSubscription(
	ctx context.Context,
	tenantID string,
) (*coreentity.InternalTenantSubscription, error) {
	return &coreentity.InternalTenantSubscription{
		ID:                "sub-1",
		Status:            coreentity.InternalTenantSubscriptionStatusActive,
		InternalProductID: "product-1",
	}, nil
}

func (billingRepoStub) UpsertSubscription(
	ctx context.Context,
	data coreentity.InternalTenantSubscription,
) (*coreentity.InternalTenantSubscription, error) {
	return &data, nil
}

func (billingRepoStub) CreateSubscriptionChange(
	ctx context.Context,
	data coreentity.InternalTenantSubscriptionChange,
) error {
	return nil
}

func (billingRepoStub) CreateEntitlementSnapshot(
	ctx context.Context,
	data coreentity.InternalEntitlementSnapshot,
) error {
	return nil
}

func (billingRepoStub) CreateLedgerEntry(
	ctx context.Context,
	data coreentity.InternalBillingLedgerEntry,
) error {
	return nil
}

func (billingRepoStub) UpdatePaymentAttemptStatus(
	ctx context.Context,
	id string,
	status string,
	metadata map[string]any,
) error {
	return nil
}

func (billingRepoStub) UpdateInvoiceStatus(
	ctx context.Context,
	id string,
	status string,
) error {
	return nil
}

func (billingRepoStub) GetInvoiceByNumber(
	ctx context.Context,
	number string,
) (*coreentity.InternalTenantInvoice, error) {
	return &coreentity.InternalTenantInvoice{
		ID:            "invoice-1",
		InvoiceNumber: number,
		Status:        coreentity.InvoiceStatusOpen,
		TenantID:      "tenant-1",
		Amount:        "999000",
	}, nil
}

func (billingRepoStub) GetPaymentAttemptByProviderRef(
	ctx context.Context,
	providerReference string,
) (*coreentity.InternalTenantPaymentAttempt, error) {
	return &coreentity.InternalTenantPaymentAttempt{
		ID:                "attempt-1",
		Status:            coreentity.PaymentAttemptStatusPending,
		ProviderReference: providerReference,
	}, nil
}

func (billingRepoStub) GetInvoice(
	ctx context.Context, tenantID, id string,
) (*coreentity.TenantBillingInvoiceDetail, error) {
	return &coreentity.TenantBillingInvoiceDetail{
		ID:             id,
		InvoiceNumber:  "TINV-001",
		Status:         coreentity.InvoiceStatusOpen,
		Amount:         "999000",
		CurrencyCode:   "IDR",
		PaymentAttempts: []coreentity.TenantBillingPaymentAttemptView{},
	}, nil
}

func (billingRepoStub) GetPaymentAttemptsByInvoice(
	ctx context.Context, tenantID, invoiceID string,
) ([]coreentity.TenantBillingPaymentAttemptView, error) {
	return nil, nil
}

func (billingRepoStub) GetInvoiceRaw(
	ctx context.Context, tenantID, id string,
) (*coreentity.InternalTenantInvoice, error) {
	return &coreentity.InternalTenantInvoice{
		ID:            id,
		TenantID:      tenantID,
		Status:        coreentity.InvoiceStatusOpen,
		InvoiceNumber: "TINV-001",
	}, nil
}

func (billingRepoStub) GetPaymentAttemptByID(
	ctx context.Context, tenantID, id string,
) (*coreentity.InternalTenantPaymentAttempt, error) {
	return &coreentity.InternalTenantPaymentAttempt{
		ID:     id,
		Status: coreentity.PaymentAttemptStatusFailed,
	}, nil
}

func (billingRepoStub) CancelInvoice(ctx context.Context, tenantID, id string) error {
	return nil
}

func (billingRepoStub) CancelPaymentAttempt(ctx context.Context, tenantID, id string) error {
	return nil
}

func (billingRepoStub) SetSubscriptionStatus(
	ctx context.Context, tenantID, subscriptionID, status string,
) error {
	return nil
}

func (billingRepoStub) GetAddOnsByIDs(
	ctx context.Context, addOnIDs []string,
) ([]coreentity.TenantBillingAddOn, error) {
	items := make([]coreentity.TenantBillingAddOn, 0, len(addOnIDs))
	for _, id := range addOnIDs {
		items = append(items, coreentity.TenantBillingAddOn{
			ID:           id,
			Name:         "Test AddOn",
			Amount:       "50000",
			Currency:     "IDR",
			BillingCycle: "monthly",
		})
	}
	return items, nil
}

func (billingRepoStub) GetAllInvoices(ctx context.Context, filter coreentity.InternalBillingInvoiceListFilter) ([]coreentity.InternalTenantInvoice, error) {
	return nil, nil
}
func (billingRepoStub) GetInvoiceByID(ctx context.Context, id string) (*coreentity.InternalTenantInvoice, error) {
	return nil, nil
}
func (billingRepoStub) GetLedgerEntries(ctx context.Context, tenantID string, filter coreentity.InternalBillingLedgerListFilter) ([]coreentity.InternalBillingLedgerEntry, error) {
	return nil, nil
}
func (billingRepoStub) GetReconciliationCases(ctx context.Context, filter coreentity.InternalBillingReconciliationCaseFilter) ([]coreentity.InternalBillingReconciliationCase, error) {
	return nil, nil
}
func (billingRepoStub) CreateInternalInvoice(ctx context.Context, data coreentity.InternalTenantInvoice) (*coreentity.InternalTenantInvoice, error) {
	data.ID = "internal-inv-1"
	return &data, nil
}
func (billingRepoStub) CreatePaymentReceipt(ctx context.Context, data coreentity.InternalPaymentReceipt) (*coreentity.InternalPaymentReceipt, error) {
	data.ID = "receipt-1"
	return &data, nil
}
func (billingRepoStub) GetPaymentReceipt(ctx context.Context, id string) (*coreentity.InternalPaymentReceipt, error) {
	return nil, nil
}
func (billingRepoStub) UpdatePaymentReceipt(ctx context.Context, data coreentity.InternalPaymentReceipt) (*coreentity.InternalPaymentReceipt, error) {
	return &data, nil
}
func (billingRepoStub) GetSubscriptionsPastPeriodEnd(ctx context.Context, limit int) ([]coreentity.InternalTenantSubscription, error) {
	return nil, nil
}
func (billingRepoStub) GetSubscriptionsInGracePastDue(ctx context.Context, limit int) ([]coreentity.InternalTenantSubscription, error) {
	return nil, nil
}
func (billingRepoStub) GetPricingInfo(ctx context.Context, pricingID string) (*coreentity.PricingInfo, error) {
	return &coreentity.PricingInfo{Amount: "100000", CurrencyCode: "IDR", BillingCycle: "monthly"}, nil
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
