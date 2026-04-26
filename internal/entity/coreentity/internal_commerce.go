package coreentity

import (
	"codebase-app/internal/entity/common"

	"github.com/shopspring/decimal"
)

const (
	QuotationStatusDraft     = "draft"
	QuotationStatusSent      = "sent"
	QuotationStatusApproved  = "approved"
	QuotationStatusRejected  = "rejected"
	QuotationStatusExpired   = "expired"
	QuotationStatusConverted = "converted"

	OrderStatusDraft                = "draft"
	OrderStatusCreatedFromQuotation = "created_from_quotation"
	OrderStatusPendingInvoice       = "pending_invoice"
	OrderStatusPendingPayment       = "pending_payment"
	OrderStatusPaid                 = "paid"
	OrderStatusCancelled            = "cancelled"
	OrderStatusExpired              = "expired"

	InvoiceStatusDraft         = "draft"
	InvoiceStatusOpen          = "open"
	InvoiceStatusPartiallyPaid = "partially_paid"
	InvoiceStatusPaid          = "paid"
	InvoiceStatusExpired       = "expired"
	InvoiceStatusCancelled     = "cancelled"

	PaymentAttemptStatusInitiated = "initiated"
	PaymentAttemptStatusPending   = "pending"
	PaymentAttemptStatusSucceeded = "succeeded"
	PaymentAttemptStatusFailed    = "failed"
	PaymentAttemptStatusExpired   = "expired"
	PaymentAttemptStatusCancelled = "cancelled"

	PaymentReceiptStatusPendingVerification = "pending_verification"
	PaymentReceiptStatusAccepted            = "accepted"
	PaymentReceiptStatusRejected            = "rejected"

	PaymentProviderDOKU   = "doku"
	PaymentProviderManual = "manual"

	OrderSourceTypeStandardPricing = "standard_pricing"
	OrderSourceTypeQuotation       = "quotation"

	ActionApproveQuotation           = "approve_quotation"
	ActionConvertQuotation           = "convert_quotation"
	ActionCreateDOKUPaymentAttempt   = "create_doku_payment_attempt"
	ActionCreateCustomPaymentAttempt = "create_custom_payment_attempt"
	ActionRetryPayment               = "retry_payment"
	ActionAcceptPaymentReceipt       = "accept_payment_receipt"
	ActionRejectPaymentReceipt       = "reject_payment_receipt"
	ActionCancelInvoice              = "cancel_invoice"
	ActionRedirectToPaymentURL       = "redirect_to_payment_url"
	ActionWaitForVerification        = "wait_for_verification"
	ActionViewInvoice                = "view_invoice"
)

type InternalCommerceAction struct {
	Key                  string
	Label                string
	Method               string
	Href                 string
	RequiresConfirmation bool
}

type InternalQuotation struct {
	UserCtx                  common.UserContext
	ID                       string
	TenantID                 string
	CompanyID                *string
	QuotationNumber          string
	Status                   string
	CurrencyCode             string
	SubtotalAmount           decimal.Decimal
	DiscountAmount           decimal.Decimal
	TaxAmount                decimal.Decimal
	TotalAmount              decimal.Decimal
	InternalProductID        string
	InternalProductPricingID string
	PricingSnapshot          map[string]any
	QuoteSnapshot            map[string]any
	ExpiresAt                *string
	ApprovedAt               *string
	ConvertedToOrderID       *string
	Metadata                 map[string]any
	CreatedAt                string
	UpdatedAt                string
}

type InternalCommerceListFilter struct {
	TenantID string
	Page     int
	Paginate int
}

type InternalOrder struct {
	UserCtx                  common.UserContext
	ID                       string
	TenantID                 string
	CompanyID                *string
	OrderNumber              string
	Status                   string
	CurrencyCode             string
	SubtotalAmount           decimal.Decimal
	DiscountAmount           decimal.Decimal
	TaxAmount                decimal.Decimal
	TotalAmount              decimal.Decimal
	InternalProductID        string
	InternalProductPricingID string
	PricingSnapshot          map[string]any
	SourceType               string
	SourceReferenceID        *string
	Metadata                 map[string]any
	OrderedAt                string
	CreatedAt                string
	UpdatedAt                string
}

type InternalInvoice struct {
	ID                string
	InternalOrderID   string
	InvoiceNumber     string
	Status            string
	CurrencyCode      string
	Amount            decimal.Decimal
	AmountPaid        decimal.Decimal
	AmountOutstanding decimal.Decimal
	DueAt             *string
	PaidAt            *string
	ExpiredAt         *string
	Metadata          map[string]any
	CreatedAt         string
	UpdatedAt         string
}

type InternalPaymentAttempt struct {
	ID                      string
	InternalInvoiceID       string
	Provider                string
	PaymentMethodType       string
	PaymentChannelCode      string
	ProviderReference       string
	ProviderRequestID       string
	ProviderPaymentURL      string
	ProviderPayloadSnapshot map[string]any
	Status                  string
	RequestedAmount         decimal.Decimal
	PaidAmount              decimal.Decimal
	ExpiredAt               *string
	PaidAt                  *string
	FailedAt                *string
	RawLastStatus           string
	Metadata                map[string]any
	CreatedAt               string
	UpdatedAt               string
}

type InternalPaymentReceipt struct {
	UserCtx                  common.UserContext
	ID                       string
	InternalInvoiceID        string
	InternalPaymentAttemptID *string
	ReceiptNumber            string
	Status                   string
	AmountReceived           decimal.Decimal
	CurrencyCode             string
	ReceivedAt               string
	VerifiedAt               *string
	VerifiedByUserID         *string
	SourceType               string
	ReferenceNumber          string
	Notes                    string
	Metadata                 map[string]any
	CreatedAt                string
	UpdatedAt                string
}

type CreateInternalQuotationInput struct {
	UserCtx                  common.UserContext
	InternalProductID        string
	InternalProductPricingID string
	CurrencyCode             string
	SubtotalAmount           decimal.Decimal
	DiscountAmount           decimal.Decimal
	TaxAmount                decimal.Decimal
	TotalAmount              decimal.Decimal
	ExpiresAt                *string
	QuoteSnapshot            map[string]any
	Metadata                 map[string]any
}

type CreateInternalOrderInput struct {
	UserCtx                  common.UserContext
	InternalProductID        string
	InternalProductPricingID string
	CurrencyCode             string
	InvoiceDueAt             *string
	Metadata                 map[string]any
}

type InternalQuotationActionInput struct {
	UserCtx      common.UserContext
	ID           string
	Action       string
	ApprovalNote string
	InvoiceDueAt *string
	RequestID    string
}

type InternalInvoiceActionInput struct {
	UserCtx            common.UserContext
	ID                 string
	Action             string
	Provider           string
	PaymentMethodType  string
	PaymentChannelCode string
	CallbackURL        string
	CustomerName       string
	CustomerEmail      string
	CustomerPhone      string
	RequestID          string
}

type InternalPaymentAttemptActionInput struct {
	UserCtx            common.UserContext
	ID                 string
	Action             string
	Provider           string
	PaymentMethodType  string
	PaymentChannelCode string
	CallbackURL        string
	CustomerName       string
	CustomerEmail      string
	CustomerPhone      string
	RequestID          string
}

type InternalCommerceBundle struct {
	Quotation *InternalQuotation
	Order     *InternalOrder
	Invoice   *InternalInvoice
}

type InternalInvoiceDetail struct {
	Invoice              InternalInvoice
	Order                InternalOrder
	Quotation            *InternalQuotation
	LatestPaymentAttempt *InternalPaymentAttempt
}

type InternalPaymentMethod struct {
	Provider          string
	PaymentMethodType string
	Label             string
}

type InternalManualPaymentInstruction struct {
	Title           string
	ReferenceNumber string
	BankName        string
	AccountNumber   string
	AccountName     string
	Notes           string
}

type InternalPaymentAttemptResult struct {
	PaymentAttempt InternalPaymentAttempt
	Instruction    *InternalManualPaymentInstruction
}

type InternalPaymentReceiptResult struct {
	Receipt InternalPaymentReceipt
	Invoice InternalInvoice
	Order   InternalOrder
}
