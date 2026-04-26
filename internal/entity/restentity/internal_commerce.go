package restentity

import "codebase-app/pkg/types"

type InternalCommerceAction struct {
	Key                  string         `json:"key"`
	Label                string         `json:"label"`
	Method               string         `json:"method"`
	Href                 string         `json:"href"`
	RequiresConfirmation bool           `json:"requires_confirmation,omitempty"`
	PayloadSchema        map[string]any `json:"payload_schema,omitempty"`
}

type CreateInternalQuotationReq struct {
	TenantID                 string         `json:"tenant_id" validate:"required"`
	InternalProductID        string         `json:"internal_product_id" validate:"required"`
	InternalProductPricingID string         `json:"internal_product_pricing_id" validate:"required"`
	CurrencyCode             string         `json:"currency_code" validate:"required,len=3"`
	SubtotalAmount           string         `json:"subtotal_amount" validate:"required,numeric"`
	DiscountAmount           string         `json:"discount_amount" validate:"required,numeric"`
	TaxAmount                string         `json:"tax_amount" validate:"required,numeric"`
	TotalAmount              string         `json:"total_amount" validate:"required,numeric"`
	ExpiresAt                *string        `json:"expires_at"`
	QuoteSnapshot            map[string]any `json:"quote_snapshot"`
	Metadata                 map[string]any `json:"metadata"`
}

type GetInternalCommerceResourceReq struct {
	ID string `params:"id" validate:"required"`
}

type GetInternalCommerceResourcesReq struct {
	types.MetaQuery
}

func (r *GetInternalCommerceResourcesReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetInternalQuotationsResp struct {
	Items []InternalQuotation `json:"items"`
	Meta  types.Meta          `json:"meta"`
}

type GetInternalOrdersResp struct {
	Items []InternalOrder `json:"items"`
	Meta  types.Meta      `json:"meta"`
}

type GetInternalInvoicesResp struct {
	Items []InternalInvoice `json:"items"`
	Meta  types.Meta        `json:"meta"`
}

type InternalQuotation struct {
	ID                       string                   `json:"id"`
	QuotationNumber          string                   `json:"quotation_number"`
	Status                   string                   `json:"status"`
	CurrencyCode             string                   `json:"currency_code"`
	SubtotalAmount           string                   `json:"subtotal_amount"`
	DiscountAmount           string                   `json:"discount_amount"`
	TaxAmount                string                   `json:"tax_amount"`
	TotalAmount              string                   `json:"total_amount"`
	InternalProductID        string                   `json:"internal_product_id"`
	InternalProductPricingID string                   `json:"internal_product_pricing_id"`
	PricingSnapshot          map[string]any           `json:"pricing_snapshot"`
	QuoteSnapshot            map[string]any           `json:"quote_snapshot"`
	ExpiresAt                *string                  `json:"expires_at"`
	ApprovedAt               *string                  `json:"approved_at"`
	ConvertedToOrderID       *string                  `json:"converted_to_order_id"`
	Metadata                 map[string]any           `json:"metadata"`
	CreatedAt                string                   `json:"created_at"`
	UpdatedAt                string                   `json:"updated_at"`
	AvailableActions         []InternalCommerceAction `json:"available_actions,omitempty"`
}

type InternalOrder struct {
	ID                       string                   `json:"id"`
	OrderNumber              string                   `json:"order_number"`
	Status                   string                   `json:"status"`
	CurrencyCode             string                   `json:"currency_code"`
	SubtotalAmount           string                   `json:"subtotal_amount"`
	DiscountAmount           string                   `json:"discount_amount"`
	TaxAmount                string                   `json:"tax_amount"`
	TotalAmount              string                   `json:"total_amount"`
	InternalProductID        string                   `json:"internal_product_id"`
	InternalProductPricingID string                   `json:"internal_product_pricing_id"`
	PricingSnapshot          map[string]any           `json:"pricing_snapshot"`
	SourceType               string                   `json:"source_type"`
	SourceReferenceID        *string                  `json:"source_reference_id"`
	Metadata                 map[string]any           `json:"metadata"`
	OrderedAt                string                   `json:"ordered_at"`
	CreatedAt                string                   `json:"created_at"`
	UpdatedAt                string                   `json:"updated_at"`
	Invoice                  *InternalInvoiceSummary  `json:"invoice,omitempty"`
	AvailableActions         []InternalCommerceAction `json:"available_actions,omitempty"`
}

type InternalInvoiceSummary struct {
	ID            string `json:"id"`
	InvoiceNumber string `json:"invoice_number"`
	Status        string `json:"status"`
}

type InternalQuotationSummary struct {
	ID              string `json:"id"`
	QuotationNumber string `json:"quotation_number"`
	Status          string `json:"status"`
}

type InternalInvoice struct {
	ID                      string                    `json:"id"`
	InternalOrderID         string                    `json:"internal_order_id"`
	InvoiceNumber           string                    `json:"invoice_number"`
	Status                  string                    `json:"status"`
	CurrencyCode            string                    `json:"currency_code"`
	Amount                  string                    `json:"amount"`
	AmountPaid              string                    `json:"amount_paid"`
	AmountOutstanding       string                    `json:"amount_outstanding"`
	DueAt                   *string                   `json:"due_at"`
	PaidAt                  *string                   `json:"paid_at"`
	ExpiredAt               *string                   `json:"expired_at"`
	Metadata                map[string]any            `json:"metadata"`
	CreatedAt               string                    `json:"created_at"`
	UpdatedAt               string                    `json:"updated_at"`
	Order                   *InternalOrder            `json:"order,omitempty"`
	Quotation               *InternalQuotationSummary `json:"quotation,omitempty"`
	LatestPaymentAttempt    *InternalPaymentAttempt   `json:"latest_payment_attempt,omitempty"`
	AvailablePaymentMethods []InternalPaymentMethod   `json:"available_payment_methods,omitempty"`
	AvailableActions        []InternalCommerceAction  `json:"available_actions,omitempty"`
}

type InternalPaymentAttempt struct {
	ID                 string                    `json:"id"`
	InternalInvoiceID  string                    `json:"internal_invoice_id"`
	Provider           string                    `json:"provider"`
	PaymentMethodType  string                    `json:"payment_method_type"`
	PaymentChannelCode string                    `json:"payment_channel_code"`
	ProviderReference  string                    `json:"provider_reference"`
	ProviderRequestID  string                    `json:"provider_request_id"`
	PaymentURL         string                    `json:"payment_url"`
	Status             string                    `json:"status"`
	RequestedAmount    string                    `json:"requested_amount"`
	PaidAmount         string                    `json:"paid_amount"`
	ExpiredAt          *string                   `json:"expired_at"`
	PaidAt             *string                   `json:"paid_at"`
	FailedAt           *string                   `json:"failed_at"`
	RawLastStatus      string                    `json:"raw_last_status"`
	Metadata           map[string]any            `json:"metadata"`
	CreatedAt          string                    `json:"created_at"`
	UpdatedAt          string                    `json:"updated_at"`
	Instruction        *ManualPaymentInstruction `json:"instruction,omitempty"`
	AvailableActions   []InternalCommerceAction  `json:"available_actions,omitempty"`
}

type InternalPaymentMethod struct {
	Provider          string `json:"provider"`
	PaymentMethodType string `json:"payment_method_type"`
	Label             string `json:"label"`
}

type ManualPaymentInstruction struct {
	Title           string `json:"title"`
	ReferenceNumber string `json:"reference_number"`
	BankName        string `json:"bank_name"`
	AccountNumber   string `json:"account_number"`
	AccountName     string `json:"account_name"`
	Notes           string `json:"notes"`
}

type ExecuteQuotationActionReq struct {
	ID           string  `params:"id" validate:"required"`
	Action       string  `json:"action" validate:"required"`
	ApprovalNote string  `json:"approval_note"`
	InvoiceDueAt *string `json:"invoice_due_at"`
}

type CreateInternalOrderReq struct {
	TenantID                 string         `json:"tenant_id" validate:"required"`
	InternalProductID        string         `json:"internal_product_id" validate:"required"`
	InternalProductPricingID string         `json:"internal_product_pricing_id" validate:"required"`
	CurrencyCode             string         `json:"currency_code" validate:"required,len=3"`
	InvoiceDueAt             *string        `json:"invoice_due_at"`
	Metadata                 map[string]any `json:"metadata"`
}

type InternalCommerceBundleResp struct {
	Quotation        *InternalQuotation       `json:"quotation,omitempty"`
	Order            *InternalOrder           `json:"order,omitempty"`
	Invoice          *InternalInvoice         `json:"invoice,omitempty"`
	AvailableActions []InternalCommerceAction `json:"available_actions,omitempty"`
}

type ExecuteInvoiceActionReq struct {
	ID                 string                      `params:"id" validate:"required"`
	Action             string                      `json:"action" validate:"required"`
	Provider           string                      `json:"provider"`
	PaymentMethodType  string                      `json:"payment_method_type"`
	PaymentChannelCode string                      `json:"payment_channel_code"`
	CallbackURL        string                      `json:"callback_url"`
	Customer           InternalCommerceCustomerReq `json:"customer"`
}

type ExecutePaymentAttemptActionReq struct {
	ID                 string                      `params:"id" validate:"required"`
	Action             string                      `json:"action" validate:"required"`
	Provider           string                      `json:"provider"`
	PaymentMethodType  string                      `json:"payment_method_type"`
	PaymentChannelCode string                      `json:"payment_channel_code"`
	CallbackURL        string                      `json:"callback_url"`
	Customer           InternalCommerceCustomerReq `json:"customer"`
}

type InternalCommerceCustomerReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type GetPaymentAttemptsResp struct {
	Items            []InternalPaymentAttempt `json:"items"`
	AvailableActions []InternalCommerceAction `json:"available_actions,omitempty"`
}

type CreatePaymentReceiptReq struct {
	InternalInvoiceID        string         `json:"internal_invoice_id" validate:"required"`
	InternalPaymentAttemptID *string        `json:"internal_payment_attempt_id"`
	Status                   string         `json:"status" validate:"required,oneof=accepted rejected pending_verification"`
	AmountReceived           string         `json:"amount_received" validate:"required,numeric"`
	CurrencyCode             string         `json:"currency_code" validate:"required,len=3"`
	ReceivedAt               string         `json:"received_at" validate:"required"`
	ReferenceNumber          string         `json:"reference_number"`
	Notes                    string         `json:"notes"`
	Metadata                 map[string]any `json:"metadata"`
}

type CreatePaymentReceiptResp struct {
	ID            string                 `json:"id"`
	ReceiptNumber string                 `json:"receipt_number"`
	Status        string                 `json:"status"`
	Invoice       InternalInvoiceSummary `json:"invoice"`
	Order         InternalOrder          `json:"order"`
}
