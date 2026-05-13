package coreentity

import (
	"encoding/json"

	"codebase-app/internal/entity/common"
)

const (
	InternalTenantSubscriptionStatusFree              = "free"
	InternalTenantSubscriptionStatusPendingActivation = "pending_activation"
	InternalTenantSubscriptionStatusActive            = "active"
	InternalTenantSubscriptionStatusGracePeriod       = "grace_period"
	InternalTenantSubscriptionStatusSuspended         = "suspended"
	InternalTenantSubscriptionStatusCancelled         = "cancelled"
	InternalTenantSubscriptionStatusExpired           = "expired"

	ResourceTypeEmployees = "employees"
	ResourceTypeBranches  = "branches"
	ResourceTypeUsers     = "users"

	InternalTenantSubscriptionTriggerInitialize                    = "initialize"
	InternalTenantSubscriptionTriggerCheckoutCreated               = "checkout_created"
	InternalTenantSubscriptionTriggerInvoicePaidFinal              = "invoice_paid_final"
	InternalTenantSubscriptionTriggerInvoiceExpiredOrAbandoned     = "invoice_expired_or_abandoned"
	InternalTenantSubscriptionTriggerRenewalUnpaid                 = "renewal_unpaid"
	InternalTenantSubscriptionTriggerRenewalPaid                   = "renewal_paid"
	InternalTenantSubscriptionTriggerGraceEnded                    = "grace_ended"
	InternalTenantSubscriptionTriggerCancelAtPeriodEnd             = "cancel_at_period_end"
	InternalTenantSubscriptionTriggerPeriodEnded                   = "period_ended"
	InternalTenantSubscriptionTriggerOutstandingPaidAndReactivated = "outstanding_paid_and_reactivated"
	InternalTenantSubscriptionTriggerAccountClosedWithoutRecovery  = "account_closed_without_recovery"
	InternalTenantSubscriptionTriggerAddOnAdded                     = "add_on_added"
	InternalTenantSubscriptionTriggerAddOnRemoved                   = "add_on_removed"

	InternalTenantBillingAccountStatusOpen      = "open"
	InternalTenantBillingAccountStatusSuspended = "suspended"
	InternalTenantBillingAccountStatusClosed    = "closed"

	InternalTenantSubscriptionChangeTypeCheckout     = "checkout"
	InternalTenantSubscriptionChangeTypeRenewal      = "renewal"
	InternalTenantSubscriptionChangeTypeUpgrade      = "upgrade"
	InternalTenantSubscriptionChangeTypeDowngrade    = "downgrade"
	InternalTenantSubscriptionChangeTypeCancellation = "cancellation"
	InternalTenantSubscriptionChangeTypeReactivation = "reactivation"
	InternalTenantSubscriptionChangeTypeExpiration   = "expiration"

	InternalBillingLedgerEntryTypeInvoiceIssued      = "invoice_issued"
	InternalBillingLedgerEntryTypePaymentSucceeded   = "payment_succeeded"
	InternalBillingLedgerEntryTypePaymentFailed      = "payment_failed"
	InternalBillingLedgerEntryTypeSubscriptionChange = "subscription_change"
	InternalBillingLedgerEntryTypeEntitlementUpdated = "entitlement_updated"
	InternalBillingLedgerEntryTypeAdjustment         = "adjustment"

	InternalBillingSourceTypeSelfServeCheckout = "self_serve_checkout"
	InternalBillingSourceTypeManualInvoice     = "manual_invoice"
	InternalBillingSourceTypeAssistedManual    = "assisted_manual"
	InternalBillingSourceTypeDOKUWebhook       = "doku_webhook"
	InternalBillingSourceTypeRenewalScheduler  = "renewal_scheduler"
	InternalBillingSourceTypeSelfServeUpgrade   = "self_serve_upgrade"

	InternalBillingReconciliationCaseStatusOpen      = "open"
	InternalBillingReconciliationCaseStatusInReview  = "in_review"
	InternalBillingReconciliationCaseStatusResolved  = "resolved"
	InternalBillingReconciliationCaseStatusDismissed = "dismissed"

	DOKUProviderStatusPending        = "Pending"
	DOKUProviderStatusSuccess        = "Success"
	DOKUProviderStatusExpired        = "Expired"
	DOKUProviderStatusCancelled      = "CANCELLED"
	DOKUProviderStatusTransactionOK  = "SUCCESS"
	DOKUProviderStatusTransactionNOK = "FAILED"
	DOKUProviderStatusOrderGenerated = "ORDER_GENERATED"
	DOKUProviderStatusOrderExpired   = "ORDER_EXPIRED"
	DOKUProviderStatusOrderRecovered = "ORDER_RECOVERED"

	DOKUProviderContextCheckoutOrder      = "checkout_order"
	DOKUProviderContextTransactionChannel = "transaction_channel"
	DOKUProviderContextOrderCheck         = "order_check"

	InternalBillingPaymentEventStatusFailedByChannel     = "failed_by_channel"
	InternalBillingPaymentEventStatusRecovered           = "recovered"
	InternalBillingPaymentEventStatusCancelledByMerchant = "cancelled_by_merchant"

	InternalEntitlementAddOnTypeFeatureUnlock       = "feature_unlock"
	InternalEntitlementAddOnTypeUsageLimitExtension = "usage_limit_extension"
)

type InternalBillingAccount struct {
	UserCtx common.UserContext

	ID       string `db:"id"`
	TenantID string `db:"tenant_id"`
	Status   string `db:"status"`
	Metadata map[string]any

	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type InternalTenantSubscription struct {
	UserCtx common.UserContext

	ID                       string
	InternalBillingAccountID string
	TenantID                 string
	ProductFamily            string
	Status                   string
	InternalProductID        string
	InternalProductPricingID string
	PlanSnapshot             map[string]any
	AddOnSnapshots           []InternalTenantSubscriptionAddOnSnapshot
	CurrentPeriodStartedAt   *string
	CurrentPeriodEndedAt     *string
	GraceEndedAt             *string
	CancelledAt              *string
	ExpiredAt                *string
	Metadata                 map[string]any

	CreatedAt string
	UpdatedAt string
}

type InternalTenantSubscriptionAddOnSnapshot struct {
	InternalProductID        string
	InternalProductPricingID string
	Type                     string
	Features                 []string
	UsageLimits              map[string]int64
	Metadata                 map[string]any
}

type InternalTenantSubscriptionChange struct {
	UserCtx common.UserContext

	ID                           string
	TenantID                     string
	InternalTenantSubscriptionID string
	ChangeType                   string
	FromStatus                   string
	ToStatus                     string
	Trigger                      string
	SourceType                   string
	SourceReferenceID            *string
	EffectiveAt                  string
	Metadata                     map[string]any

	CreatedAt string
	UpdatedAt string
}

type InternalTenantInvoice struct {
	UserCtx common.UserContext

	ID                                 string
	TenantID                           string
	InternalBillingAccountID           string
	InternalTenantSubscriptionID       *string
	InternalTenantSubscriptionChangeID *string
	InvoiceNumber                      string
	Status                             string
	CurrencyCode                       string
	Amount                             string
	AmountPaid                         string
	AmountOutstanding                  string
	SourceType                         string
	SourceReferenceID                  *string
	TargetSubscriptionState            string
	DueAt                              *string
	PaidAt                             *string
	ExpiredAt                          *string
	Metadata                           map[string]any

	CreatedAt string
	UpdatedAt string
}

type InternalTenantPaymentAttempt struct {
	UserCtx common.UserContext

	ID                      string
	TenantID                string
	InternalTenantInvoiceID string
	Provider                string
	PaymentMethodType       string
	PaymentChannelCode      string
	ProviderReference       string
	ProviderRequestID       string
	ProviderPaymentURL      string
	ProviderPayloadSnapshot map[string]any
	Status                  string
	RequestedAmount         string
	PaidAmount              string
	ExpiredAt               *string
	PaidAt                  *string
	FailedAt                *string
	CancelledAt             *string
	RawLastStatus           string
	Metadata                map[string]any

	CreatedAt string
	UpdatedAt string
}

type InternalBillingPaymentEvent struct {
	UserCtx common.UserContext

	ID                             string
	TenantID                       string
	InternalTenantInvoiceID        *string
	InternalTenantPaymentAttemptID *string
	Provider                       string
	ProviderEventID                string
	ProviderRequestID              string
	ProviderReference              string
	ProviderStatus                 string
	ProviderContext                string
	NormalizedAttemptStatus        string
	NormalizedEventStatus          string
	IsTerminal                     bool
	ProcessingResult               string
	ProcessedAt                    *string
	RawPayloadSummary              map[string]any
	Metadata                       map[string]any

	CreatedAt string
	UpdatedAt string
}

type TenantBillingPlan struct {
	ID            string
	Name          string
	Description   string
	PricingID     string
	BillingCycle  string
	Amount        string
	Currency      string
	Features      []string
	Prices        []TenantBillingPlanPrice
	AddOns        []TenantBillingAddOn
	IsCurrentPlan bool
}

type TenantBillingPlanPrice struct {
	ID                    string
	PricingID             string
	BillingCycle          string
	Amount                string
	Currency              string
	CurrencySymbol        string
	CurrencyDecimalPlaces int
	IsDefaultCurrency     bool
}

type TenantBillingAddOn struct {
	ID           string
	Code         string
	PricingID    string
	Name         string
	Description  string
	Amount       string
	Currency     string
	BillingCycle string
	RawMetadata  *json.RawMessage
}

type TenantBillingSubscriptionView struct {
	ID                 string
	Status             string
	PlanName           string
	BillingCycle       string
	CurrentPeriodStart *string
	CurrentPeriodEnd   *string
	NextRenewalAt      *string
	AddOns             []TenantBillingAddOn
}

type TenantBillingInvoice struct {
	ID               string  `db:"id"`
	InvoiceNumber    string  `db:"invoice_number"`
	Status           string  `db:"status"`
	Amount           string  `db:"amount"`
	Currency         string  `db:"currency"`
	PaymentAttemptID *string `db:"payment_attempt_id"`
	PaymentURL       *string `db:"payment_url"`
	DueAt            *string `db:"due_at"`
	PaidAt           *string `db:"paid_at"`
	CreatedAt        string  `db:"created_at"`
}

type TenantBillingInvoiceListFilter struct {
	UserCtx  common.UserContext
	Page     int
	Paginate int
}

type TenantBillingCheckoutInput struct {
	UserCtx               common.UserContext
	PriceID               string
	AddOnIDs              []string
	ReplaceActiveCheckout bool
	CallbackURL           string
	RequestID             string
}

type TenantBillingCheckoutResult struct {
	InvoiceID               string
	PaymentAttemptID        string
	PaymentURL              string
	TargetSubscriptionState string
	CheckoutReused          bool
}

type InternalEntitlementSnapshot struct {
	UserCtx common.UserContext

	ID                           string
	TenantID                     string
	InternalTenantSubscriptionID string
	SubscriptionStatus           string
	Features                     []string
	UsageLimits                  map[string]int64
	SourceSnapshot               map[string]any
	EffectiveAt                  string
	Metadata                     map[string]any

	CreatedAt string
	UpdatedAt string
}

type InternalBillingLedgerEntry struct {
	UserCtx common.UserContext

	ID                           string
	TenantID                     string
	InternalTenantSubscriptionID *string
	EntryType                    string
	SourceType                   string
	SourceReferenceID            *string
	OccurredAt                   string
	Metadata                     map[string]any

	CreatedAt string
	UpdatedAt string
}

type InternalBillingReconciliationCase struct {
	UserCtx common.UserContext

	ID                string
	TenantID          string
	Status            string
	SourceType        string
	SourceReferenceID *string
	Reason            string
	ResolutionNote    string
	Metadata          map[string]any

	CreatedAt  string
	UpdatedAt  string
	ResolvedAt *string
}

type InternalEntitlementGrant struct {
	Features    []string
	UsageLimits map[string]int64
}

type InternalEntitlementAddOn struct {
	Type        string
	Features    []string
	UsageLimits map[string]int64
}

type InternalEntitlementCalculationInput struct {
	SubscriptionStatus string
	FreeTier           InternalEntitlementGrant
	SafeAccess         InternalEntitlementGrant
	Plan               *InternalEntitlementGrant
	AddOns             []InternalEntitlementAddOn
}

type InternalBillingPaymentStatusNormalization struct {
	ProviderStatus       string
	ProviderContext      string
	PaymentAttemptStatus string
	PaymentEventStatus   string
	IsTerminal           bool
}

const (
	TenantBillingInvoiceActionCancelCheckout = "cancel_checkout_order"
)

const (
	TenantBillingPaymentAttemptActionRetry = "retry"
)

const (
	TenantBillingSubscriptionActionCancel     = "cancel"
	TenantBillingSubscriptionActionReactivate = "reactivate"
)

const (
	TenantBillingAddOnsActionAdd    = "add"
	TenantBillingAddOnsActionRemove = "remove"
)

const (
	InternalTenantInvoiceStatusPending     = "pending"
	InternalTenantInvoiceStatusPaid        = "paid"
	InternalTenantInvoiceStatusVoided      = "voided"
	InternalTenantInvoiceStatusOverdue     = "overdue"
	InternalTenantInvoiceStatusCancelled   = "cancelled"
	InternalTenantInvoiceStatusPastDue     = "past_due"
)

const (
	InternalTenantPaymentAttemptStatusInitiated  = "initiated"
	InternalTenantPaymentAttemptStatusPending    = "pending"
	InternalTenantPaymentAttemptStatusPaid       = "paid"
	InternalTenantPaymentAttemptStatusFailed     = "failed"
	InternalTenantPaymentAttemptStatusExpired    = "expired"
	InternalTenantPaymentAttemptStatusCancelled  = "cancelled"
)

type TenantBillingInvoiceDetail struct {
	ID                                 string
	InvoiceNumber                      string
	Status                             string
	CurrencyCode                       string
	Amount                             string
	AmountPaid                         string
	AmountOutstanding                  string
	SourceType                         string
	DueAt                              *string
	PaidAt                             *string
	ExpiredAt                          *string
	CreatedAt                          string
	PaymentAttempts                    []TenantBillingPaymentAttemptView
}

type TenantBillingPaymentAttemptView struct {
	ID                 string  `db:"id"`
	Provider           string  `db:"provider"`
	PaymentMethodType  string  `db:"payment_method_type"`
	PaymentChannelCode string  `db:"payment_channel_code"`
	ProviderReference  string  `db:"provider_reference"`
	PaymentURL         string  `db:"provider_payment_url"`
	Status             string  `db:"status"`
	RequestedAmount    string  `db:"requested_amount"`
	PaidAmount         string  `db:"paid_amount"`
	ExpiredAt          *string `db:"expired_at"`
	PaidAt             *string `db:"paid_at"`
	FailedAt           *string `db:"failed_at"`
	CreatedAt          string  `db:"created_at"`
}

type TenantBillingInvoiceActionInput struct {
	UserCtx common.UserContext
	ID      string
	Action  string
	Reason  string
}

type TenantBillingInvoiceActionResult struct {
	InvoiceID string
	Status    string
}

type TenantBillingPaymentAttemptActionInput struct {
	UserCtx common.UserContext
	ID      string
	Action  string
}

type TenantBillingPaymentAttemptActionResult struct {
	PaymentAttemptID string
	Status           string
	PaymentURL       string
}

type TenantBillingSubscriptionActionInput struct {
	UserCtx common.UserContext
	Action  string
}

type TenantBillingSubscriptionActionResult struct {
	SubscriptionID string
	Status         string
}

type TenantBillingAddOnsActionInput struct {
	UserCtx common.UserContext
	Action   string
	AddOnIDs []string
}

type TenantBillingAddOnsActionResult struct {
	SubscriptionID string
	AddOnIDs       []string
}

type InternalBillingInvoiceListFilter struct {
	TenantID *string
	Status   string
	FromDate string
	ToDate   string
	Page     int
	Size     int
}

type InternalBillingManualInvoiceInput struct {
	UserCtx     common.UserContext
	TenantID    string
	Amount      string
	Currency    string
	Description string
	DueAt       string
	Items       []ManualInvoiceItem
}

type ManualInvoiceItem struct {
	Description string
	Amount      string
	Quantity    int
}

type InternalBillingManualInvoiceResult struct {
	Invoice          InternalTenantInvoice
	PaymentAttemptID string
	PaymentURL       *string
}

type InternalBillingPaymentReceiptActionResult struct {
	Receipt InternalPaymentReceipt
	Status  string
}

type InternalBillingPaymentReceiptActionInput struct {
	UserCtx common.UserContext
	ID      string
	Action  string
	Reason  string
}

type InternalBillingLedgerListFilter struct {
	EntryTypes []string
	FromDate   string
	ToDate     string
	Page       int
	Size       int
}

type InternalBillingReconciliationCaseFilter struct {
	Status string
	Page   int
	Size   int
}

type PricingInfo struct {
	Amount       string
	CurrencyCode string
	BillingCycle string
}
