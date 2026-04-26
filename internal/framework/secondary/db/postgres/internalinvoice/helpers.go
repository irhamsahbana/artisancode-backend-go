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

func (r *internalInvoiceRepo) exec(ctx context.Context) transaction.SQLExecutor {
	return transaction.ExecutorFromContext(ctx, r.db)
}

func (r *internalInvoiceRepo) nextNumber(ctx context.Context, prefix string) (string, error) {
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

type invoiceRow struct {
	ID                string          `db:"id"`
	InternalOrderID   string          `db:"internal_order_id"`
	InvoiceNumber     string          `db:"invoice_number"`
	Status            string          `db:"status"`
	CurrencyCode      string          `db:"currency_code"`
	Amount            decimal.Decimal `db:"amount"`
	AmountPaid        decimal.Decimal `db:"amount_paid"`
	AmountOutstanding decimal.Decimal `db:"amount_outstanding"`
	DueAt             *string         `db:"due_at"`
	PaidAt            *string         `db:"paid_at"`
	ExpiredAt         *string         `db:"expired_at"`
	Metadata          json.RawMessage `db:"metadata"`
	CreatedAt         string          `db:"created_at"`
	UpdatedAt         *string         `db:"updated_at"`
}

type paymentAttemptRow struct {
	ID                      string          `db:"id"`
	InternalInvoiceID       string          `db:"internal_invoice_id"`
	Provider                string          `db:"provider"`
	PaymentMethodType       string          `db:"payment_method_type"`
	PaymentChannelCode      string          `db:"payment_channel_code"`
	ProviderReference       string          `db:"provider_reference"`
	ProviderRequestID       string          `db:"provider_request_id"`
	ProviderPaymentURL      string          `db:"provider_payment_url"`
	ProviderPayloadSnapshot json.RawMessage `db:"provider_payload_snapshot"`
	Status                  string          `db:"status"`
	RequestedAmount         decimal.Decimal `db:"requested_amount"`
	PaidAmount              decimal.Decimal `db:"paid_amount"`
	ExpiredAt               *string         `db:"expired_at"`
	PaidAt                  *string         `db:"paid_at"`
	FailedAt                *string         `db:"failed_at"`
	RawLastStatus           string          `db:"raw_last_status"`
	Metadata                json.RawMessage `db:"metadata"`
	CreatedAt               string          `db:"created_at"`
	UpdatedAt               *string         `db:"updated_at"`
}

func mapInvoice(row invoiceRow) *coreentity.InternalInvoice {
	item := &coreentity.InternalInvoice{
		ID:                row.ID,
		InternalOrderID:   row.InternalOrderID,
		InvoiceNumber:     row.InvoiceNumber,
		Status:            row.Status,
		CurrencyCode:      row.CurrencyCode,
		Amount:            row.Amount,
		AmountPaid:        row.AmountPaid,
		AmountOutstanding: row.AmountOutstanding,
		DueAt:             row.DueAt,
		PaidAt:            row.PaidAt,
		ExpiredAt:         row.ExpiredAt,
		Metadata:          mapJSON(row.Metadata),
		CreatedAt:         row.CreatedAt,
	}
	if row.UpdatedAt != nil {
		item.UpdatedAt = *row.UpdatedAt
	}
	return item
}

func mapPaymentAttempt(row paymentAttemptRow) *coreentity.InternalPaymentAttempt {
	item := &coreentity.InternalPaymentAttempt{
		ID:                      row.ID,
		InternalInvoiceID:       row.InternalInvoiceID,
		Provider:                row.Provider,
		PaymentMethodType:       row.PaymentMethodType,
		PaymentChannelCode:      row.PaymentChannelCode,
		ProviderReference:       row.ProviderReference,
		ProviderRequestID:       row.ProviderRequestID,
		ProviderPaymentURL:      row.ProviderPaymentURL,
		ProviderPayloadSnapshot: mapJSON(row.ProviderPayloadSnapshot),
		Status:                  row.Status,
		RequestedAmount:         row.RequestedAmount,
		PaidAmount:              row.PaidAmount,
		ExpiredAt:               row.ExpiredAt,
		PaidAt:                  row.PaidAt,
		FailedAt:                row.FailedAt,
		RawLastStatus:           row.RawLastStatus,
		Metadata:                mapJSON(row.Metadata),
		CreatedAt:               row.CreatedAt,
	}
	if row.UpdatedAt != nil {
		item.UpdatedAt = *row.UpdatedAt
	}
	return item
}
