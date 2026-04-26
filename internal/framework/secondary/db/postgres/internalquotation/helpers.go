package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/framework/secondary/db/postgres/transaction"

	"github.com/shopspring/decimal"
)

func (r *internalQuotationRepo) exec(ctx context.Context) transaction.SQLExecutor {
	return transaction.ExecutorFromContext(ctx, r.db)
}

func (r *internalQuotationRepo) nextNumber(ctx context.Context, prefix string) (string, error) {
	exec := r.exec(ctx)
	today := time.Now().Format("20060102")
	like := fmt.Sprintf("%s-%s-%%", prefix, today)
	table := map[string]string{
		"QUO": "internal_quotations",
		"ORD": "internal_orders",
		"INV": "internal_invoices",
		"RCP": "internal_payment_receipts",
	}[prefix]
	column := map[string]string{
		"QUO": "quotation_number",
		"ORD": "order_number",
		"INV": "invoice_number",
		"RCP": "receipt_number",
	}[prefix]
	query := fmt.Sprintf(`SELECT COUNT(*) + 1 FROM %s WHERE %s LIKE ?`, table, column)
	var seq int
	if err := exec.GetContext(ctx, &seq, exec.Rebind(query), like); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s-%04d", prefix, today, seq), nil
}

func normalizeCommerceListFilter(filter coreentity.InternalCommerceListFilter) coreentity.InternalCommerceListFilter {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Paginate < 1 {
		filter.Paginate = 15
	}
	if filter.Paginate > 100 {
		filter.Paginate = 100
	}
	return filter
}

func mapJSON(raw json.RawMessage) map[string]any {
	if len(raw) == 0 {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		return map[string]any{}
	}
	return out
}

type quotationRow struct {
	ID                       string          `db:"id"`
	TenantID                 string          `db:"tenant_id"`
	CompanyID                *string         `db:"company_id"`
	QuotationNumber          string          `db:"quotation_number"`
	Status                   string          `db:"status"`
	CurrencyCode             string          `db:"currency_code"`
	SubtotalAmount           decimal.Decimal `db:"subtotal_amount"`
	DiscountAmount           decimal.Decimal `db:"discount_amount"`
	TaxAmount                decimal.Decimal `db:"tax_amount"`
	TotalAmount              decimal.Decimal `db:"total_amount"`
	InternalProductID        string          `db:"internal_product_id"`
	InternalProductPricingID string          `db:"internal_product_pricing_id"`
	PricingSnapshot          json.RawMessage `db:"pricing_snapshot"`
	QuoteSnapshot            json.RawMessage `db:"quote_snapshot"`
	ExpiresAt                *string         `db:"expires_at"`
	ApprovedAt               *string         `db:"approved_at"`
	ConvertedToOrderID       *string         `db:"converted_to_order_id"`
	Metadata                 json.RawMessage `db:"metadata"`
	CreatedAt                string          `db:"created_at"`
	UpdatedAt                *string         `db:"updated_at"`
}

func mapQuotation(row quotationRow) *coreentity.InternalQuotation {
	item := &coreentity.InternalQuotation{
		ID:                       row.ID,
		TenantID:                 row.TenantID,
		CompanyID:                row.CompanyID,
		QuotationNumber:          row.QuotationNumber,
		Status:                   row.Status,
		CurrencyCode:             row.CurrencyCode,
		SubtotalAmount:           row.SubtotalAmount,
		DiscountAmount:           row.DiscountAmount,
		TaxAmount:                row.TaxAmount,
		TotalAmount:              row.TotalAmount,
		InternalProductID:        row.InternalProductID,
		InternalProductPricingID: row.InternalProductPricingID,
		PricingSnapshot:          mapJSON(row.PricingSnapshot),
		QuoteSnapshot:            mapJSON(row.QuoteSnapshot),
		ExpiresAt:                row.ExpiresAt,
		ApprovedAt:               row.ApprovedAt,
		ConvertedToOrderID:       row.ConvertedToOrderID,
		Metadata:                 mapJSON(row.Metadata),
		CreatedAt:                row.CreatedAt,
	}
	if row.UpdatedAt != nil {
		item.UpdatedAt = *row.UpdatedAt
	}
	return item
}

func mapOrder(
	id string,
	tenantID string,
	companyID *string,
	number string,
	status string,
	currency string,
	subtotal decimal.Decimal,
	discount decimal.Decimal,
	tax decimal.Decimal,
	total decimal.Decimal,
	productID string,
	pricingID string,
	pricingSnapshot json.RawMessage,
	sourceType string,
	sourceReferenceID *string,
	metadata json.RawMessage,
	orderedAt string,
	createdAt string,
	updatedAt *string,
) *coreentity.InternalOrder {
	item := &coreentity.InternalOrder{
		ID:                       id,
		TenantID:                 tenantID,
		CompanyID:                companyID,
		OrderNumber:              number,
		Status:                   status,
		CurrencyCode:             currency,
		SubtotalAmount:           subtotal,
		DiscountAmount:           discount,
		TaxAmount:                tax,
		TotalAmount:              total,
		InternalProductID:        productID,
		InternalProductPricingID: pricingID,
		PricingSnapshot:          mapJSON(pricingSnapshot),
		SourceType:               sourceType,
		SourceReferenceID:        sourceReferenceID,
		Metadata:                 mapJSON(metadata),
		OrderedAt:                orderedAt,
		CreatedAt:                createdAt,
	}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}
	return item
}

func mapInvoice(
	id string,
	orderID string,
	number string,
	status string,
	currency string,
	amount decimal.Decimal,
	paid decimal.Decimal,
	outstanding decimal.Decimal,
	dueAt *string,
	paidAt *string,
	expiredAt *string,
	metadata json.RawMessage,
	createdAt string,
	updatedAt *string,
) *coreentity.InternalInvoice {
	item := &coreentity.InternalInvoice{
		ID:                id,
		InternalOrderID:   orderID,
		InvoiceNumber:     number,
		Status:            status,
		CurrencyCode:      currency,
		Amount:            amount,
		AmountPaid:        paid,
		AmountOutstanding: outstanding,
		DueAt:             dueAt,
		PaidAt:            paidAt,
		ExpiredAt:         expiredAt,
		Metadata:          mapJSON(metadata),
		CreatedAt:         createdAt,
	}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}
	return item
}

func mapPaymentAttempt(
	id string,
	invoiceID string,
	provider string,
	methodType string,
	channel string,
	providerReference string,
	requestID string,
	paymentURL string,
	payload json.RawMessage,
	status string,
	requested decimal.Decimal,
	paid decimal.Decimal,
	expiredAt *string,
	paidAt *string,
	failedAt *string,
	rawStatus string,
	metadata json.RawMessage,
	createdAt string,
	updatedAt *string,
) *coreentity.InternalPaymentAttempt {
	item := &coreentity.InternalPaymentAttempt{
		ID:                      id,
		InternalInvoiceID:       invoiceID,
		Provider:                provider,
		PaymentMethodType:       methodType,
		PaymentChannelCode:      channel,
		ProviderReference:       providerReference,
		ProviderRequestID:       requestID,
		ProviderPaymentURL:      paymentURL,
		ProviderPayloadSnapshot: mapJSON(payload),
		Status:                  status,
		RequestedAmount:         requested,
		PaidAmount:              paid,
		ExpiredAt:               expiredAt,
		PaidAt:                  paidAt,
		FailedAt:                failedAt,
		RawLastStatus:           rawStatus,
		Metadata:                mapJSON(metadata),
		CreatedAt:               createdAt,
	}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}
	return item
}
