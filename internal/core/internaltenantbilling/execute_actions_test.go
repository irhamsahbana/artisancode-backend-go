package core

import (
	"context"
	"database/sql"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/ports/integration/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type executeActionsBillingStub struct {
	invoiceStatus        string
	attemptStatus        string
	subscriptionStatus   string
	subscriptionID       string
	invoiceNotFound      bool
	attemptNotFound      bool
	subscriptionNotFound bool
	addOnsReturn         []coreentity.TenantBillingAddOn
	initialAddOnSnaps    []coreentity.InternalTenantSubscriptionAddOnSnapshot
}

func (s *executeActionsBillingStub) EnsureBillingAccount(ctx context.Context, tenantID string) (*coreentity.InternalBillingAccount, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) GetInvoices(ctx context.Context, filter coreentity.TenantBillingInvoiceListFilter) ([]coreentity.TenantBillingInvoice, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) CreateInvoice(ctx context.Context, data coreentity.InternalTenantInvoice) (*coreentity.InternalTenantInvoice, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) CreatePaymentAttempt(ctx context.Context, data coreentity.InternalTenantPaymentAttempt) (*coreentity.InternalTenantPaymentAttempt, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) UpdatePaymentAttemptGateway(ctx context.Context, data coreentity.InternalTenantPaymentAttempt) (*coreentity.InternalTenantPaymentAttempt, error) {
	return &data, nil
}
func (s *executeActionsBillingStub) GetSubscription(ctx context.Context, tenantID string) (*coreentity.InternalTenantSubscription, error) {
	if s.subscriptionNotFound {
		return nil, sql.ErrNoRows
	}
	return &coreentity.InternalTenantSubscription{
		ID:       s.subscriptionID,
		TenantID: tenantID,
		Status:   s.subscriptionStatus,
	}, nil
}
func (s *executeActionsBillingStub) GetLatestEntitlementSnapshot(ctx context.Context, tenantID string) (*coreentity.InternalEntitlementSnapshot, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) GetActiveSubscription(ctx context.Context, tenantID string) (*coreentity.InternalTenantSubscription, error) {
	if s.subscriptionNotFound {
		return nil, sql.ErrNoRows
	}
	return &coreentity.InternalTenantSubscription{
		ID:              s.subscriptionID,
		TenantID:        tenantID,
		Status:          s.subscriptionStatus,
		AddOnSnapshots:  s.initialAddOnSnaps,
	}, nil
}
func (s *executeActionsBillingStub) UpsertSubscription(ctx context.Context, data coreentity.InternalTenantSubscription) (*coreentity.InternalTenantSubscription, error) {
	return &data, nil
}
func (s *executeActionsBillingStub) CreateSubscriptionChange(ctx context.Context, data coreentity.InternalTenantSubscriptionChange) error {
	return nil
}
func (s *executeActionsBillingStub) CreateEntitlementSnapshot(ctx context.Context, data coreentity.InternalEntitlementSnapshot) error {
	return nil
}
func (s *executeActionsBillingStub) CreateLedgerEntry(ctx context.Context, data coreentity.InternalBillingLedgerEntry) error {
	return nil
}
func (s *executeActionsBillingStub) UpdatePaymentAttemptStatus(ctx context.Context, id string, status string, metadata map[string]any) error {
	return nil
}
func (s *executeActionsBillingStub) UpdateInvoiceStatus(ctx context.Context, id string, status string) error {
	return nil
}
func (s *executeActionsBillingStub) GetInvoiceByNumber(ctx context.Context, number string) (*coreentity.InternalTenantInvoice, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) GetPaymentAttemptByProviderRef(ctx context.Context, providerReference string) (*coreentity.InternalTenantPaymentAttempt, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) GetInvoice(ctx context.Context, tenantID string, id string) (*coreentity.TenantBillingInvoiceDetail, error) {
	if s.invoiceNotFound {
		return nil, sql.ErrNoRows
	}
	return &coreentity.TenantBillingInvoiceDetail{
		ID:             id,
		InvoiceNumber:  "INV-001",
		Status:         s.invoiceStatus,
		Amount:         "100000",
		CurrencyCode:   "IDR",
		PaymentAttempts: []coreentity.TenantBillingPaymentAttemptView{},
	}, nil
}
func (s *executeActionsBillingStub) GetPaymentAttemptsByInvoice(ctx context.Context, tenantID string, invoiceID string) ([]coreentity.TenantBillingPaymentAttemptView, error) {
	return []coreentity.TenantBillingPaymentAttemptView{}, nil
}
func (s *executeActionsBillingStub) GetInvoiceRaw(ctx context.Context, tenantID string, id string) (*coreentity.InternalTenantInvoice, error) {
	if s.invoiceNotFound {
		return nil, sql.ErrNoRows
	}
	return &coreentity.InternalTenantInvoice{
		ID:            id,
		TenantID:      tenantID,
		Status:        s.invoiceStatus,
		InvoiceNumber: "INV-001",
		Amount:        "100000",
		CurrencyCode:  "IDR",
	}, nil
}
func (s *executeActionsBillingStub) GetPaymentAttemptByID(ctx context.Context, tenantID string, id string) (*coreentity.InternalTenantPaymentAttempt, error) {
	if s.attemptNotFound {
		return nil, sql.ErrNoRows
	}
	return &coreentity.InternalTenantPaymentAttempt{
		ID:     id,
		Status: s.attemptStatus,
	}, nil
}
func (s *executeActionsBillingStub) CancelInvoice(ctx context.Context, tenantID string, id string) error {
	return nil
}
func (s *executeActionsBillingStub) CancelPaymentAttempt(ctx context.Context, tenantID string, id string) error {
	return nil
}
func (s *executeActionsBillingStub) SetSubscriptionStatus(ctx context.Context, tenantID string, subscriptionID string, status string) error {
	return nil
}
func (s *executeActionsBillingStub) GetAddOnsByIDs(ctx context.Context, addOnIDs []string) ([]coreentity.TenantBillingAddOn, error) {
	return s.addOnsReturn, nil
}

func (s *executeActionsBillingStub) GetAllInvoices(ctx context.Context, filter coreentity.InternalBillingInvoiceListFilter) ([]coreentity.InternalTenantInvoice, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) GetInvoiceByID(ctx context.Context, id string) (*coreentity.InternalTenantInvoice, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) GetLedgerEntries(ctx context.Context, tenantID string, filter coreentity.InternalBillingLedgerListFilter) ([]coreentity.InternalBillingLedgerEntry, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) GetReconciliationCases(ctx context.Context, filter coreentity.InternalBillingReconciliationCaseFilter) ([]coreentity.InternalBillingReconciliationCase, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) CreateInternalInvoice(ctx context.Context, data coreentity.InternalTenantInvoice) (*coreentity.InternalTenantInvoice, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) CreatePaymentReceipt(ctx context.Context, data coreentity.InternalPaymentReceipt) (*coreentity.InternalPaymentReceipt, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) GetPaymentReceipt(ctx context.Context, id string) (*coreentity.InternalPaymentReceipt, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) UpdatePaymentReceipt(ctx context.Context, data coreentity.InternalPaymentReceipt) (*coreentity.InternalPaymentReceipt, error) {
	return &data, nil
}
func (s *executeActionsBillingStub) GetSubscriptionsPastPeriodEnd(ctx context.Context, limit int) ([]coreentity.InternalTenantSubscription, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) GetSubscriptionsInGracePastDue(ctx context.Context, limit int) ([]coreentity.InternalTenantSubscription, error) {
	return nil, nil
}
func (s *executeActionsBillingStub) GetPricingInfo(ctx context.Context, pricingID string) (*coreentity.PricingInfo, error) {
	return &coreentity.PricingInfo{Amount: "100000", CurrencyCode: "IDR", BillingCycle: "monthly"}, nil
}

type stubProductRepo struct{}

func (s *stubProductRepo) GetInternalProducts(ctx context.Context, filter coreentity.InternalProductListFilter) ([]coreentity.InternalProduct, int, error) {
	return nil, 0, nil
}
func (s *stubProductRepo) GetInternalProduct(ctx context.Context, filter coreentity.InternalProductFilter) (*coreentity.InternalProduct, error) {
	return &coreentity.InternalProduct{ID: filter.ID, Code: "CORE_HR"}, nil
}
func (s *stubProductRepo) CreateInternalProduct(ctx context.Context, data coreentity.InternalProduct) (*coreentity.InternalProduct, error) {
	return nil, nil
}
func (s *stubProductRepo) UpdateInternalProduct(ctx context.Context, data coreentity.InternalProduct) error {
	return nil
}
func (s *stubProductRepo) DeleteInternalProduct(ctx context.Context, filter coreentity.InternalProductDeleteFilter) error {
	return nil
}
func (s *stubProductRepo) ExistsInternalProductByCode(ctx context.Context, code, excludeID string) (bool, error) {
	return false, nil
}
func (s *stubProductRepo) GetInternalProductPricings(ctx context.Context, filter coreentity.InternalProductPricingListFilter) ([]coreentity.InternalProductPricing, int, error) {
	return nil, 0, nil
}
func (s *stubProductRepo) GetInternalProductPricing(ctx context.Context, filter coreentity.InternalProductPricingFilter) (*coreentity.InternalProductPricing, error) {
	return nil, nil
}
func (s *stubProductRepo) CreateInternalProductPricing(ctx context.Context, data coreentity.InternalProductPricing) (*coreentity.InternalProductPricing, error) {
	return nil, nil
}
func (s *stubProductRepo) UpdateInternalProductPricing(ctx context.Context, data coreentity.InternalProductPricing) error {
	return nil
}
func (s *stubProductRepo) DeleteInternalProductPricing(ctx context.Context, filter coreentity.InternalProductPricingDeleteFilter) error {
	return nil
}
func (s *stubProductRepo) ExistsInternalProductPricingByCode(ctx context.Context, internalProductID, code, excludeID string) (bool, error) {
	return false, nil
}
func (s *stubProductRepo) GetInternalProductPrices(ctx context.Context, filter coreentity.InternalProductPriceListFilter) ([]coreentity.InternalProductPrice, int, error) {
	return nil, 0, nil
}
func (s *stubProductRepo) GetInternalProductPrice(ctx context.Context, filter coreentity.InternalProductPriceFilter) (*coreentity.InternalProductPrice, error) {
	return nil, nil
}
func (s *stubProductRepo) CreateInternalProductPrice(ctx context.Context, data coreentity.InternalProductPrice) (*coreentity.InternalProductPrice, error) {
	return nil, nil
}
func (s *stubProductRepo) UpdateInternalProductPrice(ctx context.Context, data coreentity.InternalProductPrice) error {
	return nil
}
func (s *stubProductRepo) DeleteInternalProductPrice(ctx context.Context, filter coreentity.InternalProductPriceDeleteFilter) error {
	return nil
}
func (s *stubProductRepo) ExistsOverlappingInternalProductPrice(ctx context.Context, filter coreentity.InternalProductPriceOverlapFilter) (bool, error) {
	return false, nil
}

func TestGetInvoice(t *testing.T) {
	core := NewInternalTenantBillingCore(Config{
		BillingRepo: &executeActionsBillingStub{},
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		detail, err := core.GetInvoice(ctx, "inv-1")
		require.NoError(t, err)
		require.Equal(t, "inv-1", detail.ID)
	})

	t.Run("not_found", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			BillingRepo: &executeActionsBillingStub{invoiceNotFound: true},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		_, err := core.GetInvoice(ctx, "inv-not-found")
		require.Error(t, err)
	})
}

func TestGetPaymentAttempts(t *testing.T) {
	core := NewInternalTenantBillingCore(Config{
		BillingRepo: &executeActionsBillingStub{},
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		items, err := core.GetPaymentAttempts(ctx, "inv-1")
		require.NoError(t, err)
		require.NotNil(t, items)
	})

	t.Run("invoice_not_found", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			BillingRepo: &executeActionsBillingStub{invoiceNotFound: true},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		_, err := core.GetPaymentAttempts(ctx, "inv-not-found")
		require.Error(t, err)
	})
}

func TestExecuteInvoiceAction(t *testing.T) {
	t.Run("cancel_checkout_order", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			BillingRepo: &executeActionsBillingStub{invoiceStatus: coreentity.InternalTenantInvoiceStatusPending},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		result, err := core.ExecuteInvoiceAction(ctx, coreentity.TenantBillingInvoiceActionInput{
			ID:     "inv-1",
			Action: coreentity.TenantBillingInvoiceActionCancelCheckout,
		})
		require.NoError(t, err)
		require.Equal(t, coreentity.InternalTenantInvoiceStatusCancelled, result.Status)
	})

	t.Run("invoice_not_found", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			BillingRepo: &executeActionsBillingStub{invoiceNotFound: true},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		_, err := core.ExecuteInvoiceAction(ctx, coreentity.TenantBillingInvoiceActionInput{
			ID:     "inv-not-found",
			Action: coreentity.TenantBillingInvoiceActionCancelCheckout,
		})
		require.Error(t, err)
	})

	t.Run("unsupported_action", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			BillingRepo: &executeActionsBillingStub{
				subscriptionStatus: coreentity.InternalTenantSubscriptionStatusActive,
				subscriptionID:     "sub-1",
			},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		_, err := core.ExecuteInvoiceAction(ctx, coreentity.TenantBillingInvoiceActionInput{
			ID:     "inv-1",
			Action: "invalid_action",
		})
		require.Error(t, err)
	})
}

func TestExecutePaymentAttemptAction(t *testing.T) {
	t.Run("retry_failed", func(t *testing.T) {
		doku := mocks.NewDokuClient(t)
		doku.EXPECT().CreatePayment(mock.Anything, mock.Anything).Return(&restentity.DokuCreatePaymentResponse{
			InvoiceNumber: "INV-001",
			Amount:        100000,
			PaymentURL:    "https://pay.doku.example/test",
			RequestID:     "req-1",
		}, nil)

		core := NewInternalTenantBillingCore(Config{
			BillingRepo: &executeActionsBillingStub{
				attemptStatus: coreentity.InternalTenantPaymentAttemptStatusFailed,
				invoiceStatus: coreentity.InternalTenantInvoiceStatusPending,
			},
			DOKU: doku,
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		result, err := core.ExecutePaymentAttemptAction(ctx, coreentity.TenantBillingPaymentAttemptActionInput{
			ID:     "pa-1",
			Action: coreentity.TenantBillingPaymentAttemptActionRetry,
		})
		require.NoError(t, err)
		require.Equal(t, coreentity.InternalTenantPaymentAttemptStatusPending, result.Status)
	})

	t.Run("retry_expired", func(t *testing.T) {
		doku := mocks.NewDokuClient(t)
		doku.EXPECT().CreatePayment(mock.Anything, mock.Anything).Return(&restentity.DokuCreatePaymentResponse{
			InvoiceNumber: "INV-001",
			Amount:        100000,
			PaymentURL:    "https://pay.doku.example/test",
			RequestID:     "req-1",
		}, nil)

		core := NewInternalTenantBillingCore(Config{
			BillingRepo: &executeActionsBillingStub{
				attemptStatus: coreentity.InternalTenantPaymentAttemptStatusExpired,
				invoiceStatus: coreentity.InternalTenantInvoiceStatusPending,
			},
			DOKU: doku,
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		result, err := core.ExecutePaymentAttemptAction(ctx, coreentity.TenantBillingPaymentAttemptActionInput{
			ID:     "pa-expired",
			Action: coreentity.TenantBillingPaymentAttemptActionRetry,
		})
		require.NoError(t, err)
		require.Equal(t, coreentity.InternalTenantPaymentAttemptStatusPending, result.Status)
	})

	t.Run("not_found", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			BillingRepo: &executeActionsBillingStub{attemptNotFound: true},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		_, err := core.ExecutePaymentAttemptAction(ctx, coreentity.TenantBillingPaymentAttemptActionInput{
			ID:     "pa-not-found",
			Action: coreentity.TenantBillingPaymentAttemptActionRetry,
		})
		require.Error(t, err)
	})

	t.Run("unsupported_action", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			ProductRepo: &stubProductRepo{},
			BillingRepo: &executeActionsBillingStub{},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		_, err := core.ExecutePaymentAttemptAction(ctx, coreentity.TenantBillingPaymentAttemptActionInput{
			ID:     "pa-1",
			Action: "invalid_action",
		})
		require.Error(t, err)
	})
}

func TestExecuteSubscriptionAction(t *testing.T) {
	t.Run("cancel_active", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			BillingRepo: &executeActionsBillingStub{
				subscriptionStatus: coreentity.InternalTenantSubscriptionStatusActive,
				subscriptionID:     "sub-1",
			},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		result, err := core.ExecuteSubscriptionAction(ctx, coreentity.TenantBillingSubscriptionActionInput{
			Action: coreentity.TenantBillingSubscriptionActionCancel,
		})
		require.NoError(t, err)
		require.Equal(t, coreentity.InternalTenantSubscriptionStatusCancelled, result.Status)
	})

	t.Run("reactivate_suspended", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			BillingRepo: &executeActionsBillingStub{
				subscriptionStatus: coreentity.InternalTenantSubscriptionStatusSuspended,
				subscriptionID:     "sub-2",
			},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		result, err := core.ExecuteSubscriptionAction(ctx, coreentity.TenantBillingSubscriptionActionInput{
			Action: coreentity.TenantBillingSubscriptionActionReactivate,
		})
		require.NoError(t, err)
		require.Equal(t, coreentity.InternalTenantSubscriptionStatusActive, result.Status)
	})

	t.Run("cancel_not_active", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			BillingRepo: &executeActionsBillingStub{
				subscriptionStatus: coreentity.InternalTenantSubscriptionStatusFree,
				subscriptionID:     "sub-3",
			},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		_, err := core.ExecuteSubscriptionAction(ctx, coreentity.TenantBillingSubscriptionActionInput{
			Action: coreentity.TenantBillingSubscriptionActionCancel,
		})
		require.Error(t, err)
	})

	t.Run("unsupported_action", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			BillingRepo: &executeActionsBillingStub{
				subscriptionStatus: coreentity.InternalTenantSubscriptionStatusActive,
				subscriptionID:     "sub-1",
			},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		_, err := core.ExecuteSubscriptionAction(ctx, coreentity.TenantBillingSubscriptionActionInput{
			Action: "invalid_action",
		})
		require.Error(t, err)
	})
}

func TestExecuteAddOnsAction(t *testing.T) {
	t.Run("add_success", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			ProductRepo: &stubProductRepo{},
			BillingRepo: &executeActionsBillingStub{
				addOnsReturn: []coreentity.TenantBillingAddOn{
					{ID: "addon-1", Name: "Test 1"},
					{ID: "addon-2", Name: "Test 2"},
				},
				subscriptionStatus: coreentity.InternalTenantSubscriptionStatusActive,
				subscriptionID:     "sub-1",
			},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		result, err := core.ExecuteAddOnsAction(ctx, coreentity.TenantBillingAddOnsActionInput{
			Action:   coreentity.TenantBillingAddOnsActionAdd,
			AddOnIDs: []string{"addon-1", "addon-2"},
		})
		require.NoError(t, err)
		require.Len(t, result.AddOnIDs, 2)
	})

	t.Run("add_partial_mismatch", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			ProductRepo: &stubProductRepo{},
			BillingRepo: &executeActionsBillingStub{
				addOnsReturn: []coreentity.TenantBillingAddOn{
					{ID: "addon-1", Name: "Test 1"},
				},
				subscriptionStatus: coreentity.InternalTenantSubscriptionStatusActive,
				subscriptionID:     "sub-1",
			},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		_, err := core.ExecuteAddOnsAction(ctx, coreentity.TenantBillingAddOnsActionInput{
			Action:   coreentity.TenantBillingAddOnsActionAdd,
			AddOnIDs: []string{"addon-1", "addon-missing"},
		})
		require.Error(t, err)
	})

	t.Run("add_empty", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			ProductRepo: &stubProductRepo{},
			BillingRepo: &executeActionsBillingStub{},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		_, err := core.ExecuteAddOnsAction(ctx, coreentity.TenantBillingAddOnsActionInput{
			Action:   coreentity.TenantBillingAddOnsActionAdd,
			AddOnIDs: []string{},
		})
		require.Error(t, err)
	})

	t.Run("remove_success", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			ProductRepo: &stubProductRepo{},
			BillingRepo: &executeActionsBillingStub{
				subscriptionStatus: coreentity.InternalTenantSubscriptionStatusActive,
				subscriptionID:     "sub-1",
				initialAddOnSnaps: []coreentity.InternalTenantSubscriptionAddOnSnapshot{
					{InternalProductID: "addon-1"},
					{InternalProductID: "addon-2"},
				},
			},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		result, err := core.ExecuteAddOnsAction(ctx, coreentity.TenantBillingAddOnsActionInput{
			Action:   coreentity.TenantBillingAddOnsActionRemove,
			AddOnIDs: []string{"addon-1"},
		})
		require.NoError(t, err)
		require.Len(t, result.AddOnIDs, 1)
	})

	t.Run("unsupported_action", func(t *testing.T) {
		core := NewInternalTenantBillingCore(Config{
			BillingRepo: &executeActionsBillingStub{},
		})
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		_, err := core.ExecuteAddOnsAction(ctx, coreentity.TenantBillingAddOnsActionInput{
			Action: "invalid_action",
		})
		require.Error(t, err)
	})
}
