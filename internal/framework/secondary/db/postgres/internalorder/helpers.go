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

func (r *internalOrderRepo) exec(ctx context.Context) transaction.SQLExecutor {
	return transaction.ExecutorFromContext(ctx, r.db)
}

func (r *internalOrderRepo) nextNumber(ctx context.Context, prefix string) (string, error) {
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

type orderRow struct {
	ID                       string          `db:"id"`
	TenantID                 string          `db:"tenant_id"`
	CompanyID                *string         `db:"company_id"`
	OrderNumber              string          `db:"order_number"`
	Status                   string          `db:"status"`
	CurrencyCode             string          `db:"currency_code"`
	SubtotalAmount           decimal.Decimal `db:"subtotal_amount"`
	DiscountAmount           decimal.Decimal `db:"discount_amount"`
	TaxAmount                decimal.Decimal `db:"tax_amount"`
	TotalAmount              decimal.Decimal `db:"total_amount"`
	InternalProductID        string          `db:"internal_product_id"`
	InternalProductPricingID string          `db:"internal_product_pricing_id"`
	PricingSnapshot          json.RawMessage `db:"pricing_snapshot"`
	SourceType               string          `db:"source_type"`
	SourceReferenceID        *string         `db:"source_reference_id"`
	Metadata                 json.RawMessage `db:"metadata"`
	OrderedAt                string          `db:"ordered_at"`
	CreatedAt                string          `db:"created_at"`
	UpdatedAt                *string         `db:"updated_at"`
}

func mapOrder(row orderRow) *coreentity.InternalOrder {
	item := &coreentity.InternalOrder{
		ID:                       row.ID,
		TenantID:                 row.TenantID,
		CompanyID:                row.CompanyID,
		OrderNumber:              row.OrderNumber,
		Status:                   row.Status,
		CurrencyCode:             row.CurrencyCode,
		SubtotalAmount:           row.SubtotalAmount,
		DiscountAmount:           row.DiscountAmount,
		TaxAmount:                row.TaxAmount,
		TotalAmount:              row.TotalAmount,
		InternalProductID:        row.InternalProductID,
		InternalProductPricingID: row.InternalProductPricingID,
		PricingSnapshot:          mapJSON(row.PricingSnapshot),
		SourceType:               row.SourceType,
		SourceReferenceID:        row.SourceReferenceID,
		Metadata:                 mapJSON(row.Metadata),
		OrderedAt:                row.OrderedAt,
		CreatedAt:                row.CreatedAt,
	}
	if row.UpdatedAt != nil {
		item.UpdatedAt = *row.UpdatedAt
	}
	return item
}
