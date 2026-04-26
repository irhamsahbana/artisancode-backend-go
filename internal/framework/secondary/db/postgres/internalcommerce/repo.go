package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/framework/secondary/db/postgres/transaction"
	"codebase-app/internal/infrastructure/tracing"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

type internalCommerceRepo struct {
	db *sqlx.DB
}

type Config struct {
	DB *sqlx.DB
}

var _ portsRepo.InternalCommerceRepository = &internalCommerceRepo{}

func NewInternalCommerceRepository(cfg Config) portsRepo.InternalCommerceRepository {
	return &internalCommerceRepo{db: cfg.DB}
}

func (r *internalCommerceRepo) exec(ctx context.Context) transaction.SQLExecutor {
	return transaction.ExecutorFromContext(ctx, r.db)
}

func (r *internalCommerceRepo) nextNumber(ctx context.Context, prefix string) (string, error) {
	exec := r.exec(ctx)
	today := time.Now().Format("20060102")
	like := fmt.Sprintf("%s-%s-%%", prefix, today)
	table := map[string]string{"QUO": "internal_quotations", "ORD": "internal_orders", "INV": "internal_invoices", "RCP": "internal_payment_receipts"}[prefix]
	column := map[string]string{"QUO": "quotation_number", "ORD": "order_number", "INV": "invoice_number", "RCP": "receipt_number"}[prefix]
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

func (r *internalCommerceRepo) GetQuotations(ctx context.Context, filter coreentity.InternalCommerceListFilter) ([]coreentity.InternalQuotation, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:GetQuotations")
	defer span.End()

	filter = normalizeCommerceListFilter(filter)
	type row struct {
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
	where := `WHERE deleted_at IS NULL`
	args := []any{}
	if filter.TenantID != "" {
		where += ` AND tenant_id = ?`
		args = append(args, filter.TenantID)
	}
	var total int
	countQuery := `SELECT COUNT(*) FROM internal_quotations ` + where
	if err := r.exec(ctx).GetContext(ctx, &total, r.exec(ctx).Rebind(countQuery), args...); err != nil {
		return nil, 0, err
	}
	rows := make([]row, 0)
	query := `
		SELECT id, tenant_id, company_id, quotation_number, status, currency_code,
			subtotal_amount, discount_amount, tax_amount, total_amount,
			internal_product_id, internal_product_pricing_id, pricing_snapshot, quote_snapshot,
			expires_at, approved_at, converted_to_order_id, metadata, created_at, updated_at
		FROM internal_quotations
		` + where + `
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)
	if err := r.exec(ctx).SelectContext(ctx, &rows, r.exec(ctx).Rebind(query), args...); err != nil {
		return nil, 0, err
	}
	items := make([]coreentity.InternalQuotation, 0, len(rows))
	for _, item := range rows {
		items = append(items, *mapQuotation(item.ID, item.TenantID, item.CompanyID, item.QuotationNumber, item.Status, item.CurrencyCode, item.SubtotalAmount, item.DiscountAmount, item.TaxAmount, item.TotalAmount, item.InternalProductID, item.InternalProductPricingID, item.PricingSnapshot, item.QuoteSnapshot, item.ExpiresAt, item.ApprovedAt, item.ConvertedToOrderID, item.Metadata, item.CreatedAt, item.UpdatedAt))
	}
	return items, total, nil
}

func (r *internalCommerceRepo) CreateQuotation(ctx context.Context, data coreentity.InternalQuotation) (*coreentity.InternalQuotation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:CreateQuotation")
	defer span.End()

	exec := r.exec(ctx)
	number, err := r.nextNumber(ctx, "QUO")
	if err != nil {
		return nil, err
	}
	pricingSnapshot, _ := json.Marshal(data.PricingSnapshot)
	quoteSnapshot, _ := json.Marshal(data.QuoteSnapshot)
	metadata, _ := json.Marshal(data.Metadata)
	query := `
		INSERT INTO internal_quotations (
			tenant_id, company_id, quotation_number, status, currency_code,
			subtotal_amount, discount_amount, tax_amount, total_amount,
			internal_product_id, internal_product_pricing_id, pricing_snapshot, quote_snapshot, expires_at, metadata
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, quotation_number, status, created_at
	`
	if err := exec.QueryRowxContext(ctx, exec.Rebind(query),
		data.TenantID, data.CompanyID, number, data.Status, data.CurrencyCode,
		data.SubtotalAmount, data.DiscountAmount, data.TaxAmount, data.TotalAmount,
		data.InternalProductID, data.InternalProductPricingID, pricingSnapshot, quoteSnapshot, data.ExpiresAt, metadata,
	).Scan(&data.ID, &data.QuotationNumber, &data.Status, &data.CreatedAt); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create internal quotation")
		return nil, err
	}
	return &data, nil
}

func (r *internalCommerceRepo) GetQuotation(ctx context.Context, tenantID, id string) (*coreentity.InternalQuotation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:GetQuotation")
	defer span.End()

	type row struct {
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
	var item row
	query := `
		SELECT id, tenant_id, company_id, quotation_number, status, currency_code,
			subtotal_amount, discount_amount, tax_amount, total_amount,
			internal_product_id, internal_product_pricing_id, pricing_snapshot, quote_snapshot,
			expires_at, approved_at, converted_to_order_id, metadata, created_at, updated_at
		FROM internal_quotations
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`
	if err := r.exec(ctx).GetContext(ctx, &item, r.exec(ctx).Rebind(query), id, tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Quotation not found")
		}
		return nil, err
	}
	return mapQuotation(item.ID, item.TenantID, item.CompanyID, item.QuotationNumber, item.Status, item.CurrencyCode, item.SubtotalAmount, item.DiscountAmount, item.TaxAmount, item.TotalAmount, item.InternalProductID, item.InternalProductPricingID, item.PricingSnapshot, item.QuoteSnapshot, item.ExpiresAt, item.ApprovedAt, item.ConvertedToOrderID, item.Metadata, item.CreatedAt, item.UpdatedAt), nil
}

func (r *internalCommerceRepo) ApproveQuotation(ctx context.Context, tenantID, id string) (*coreentity.InternalQuotation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:ApproveQuotation")
	defer span.End()

	query := `UPDATE internal_quotations SET status = 'approved', approved_at = NOW(), updated_at = NOW() WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`
	result, err := r.exec(ctx).ExecContext(ctx, r.exec(ctx).Rebind(query), id, tenantID)
	if err != nil {
		return nil, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return nil, errmsg.NewCustomErrors(404).SetMessage("Quotation not found")
	}
	return r.GetQuotation(ctx, tenantID, id)
}

func (r *internalCommerceRepo) MarkQuotationConverted(ctx context.Context, tenantID, quotationID, orderID string) (*coreentity.InternalQuotation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:MarkQuotationConverted")
	defer span.End()

	query := `UPDATE internal_quotations SET status = 'converted', converted_to_order_id = ?, updated_at = NOW() WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`
	result, err := r.exec(ctx).ExecContext(ctx, r.exec(ctx).Rebind(query), orderID, quotationID, tenantID)
	if err != nil {
		return nil, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return nil, errmsg.NewCustomErrors(404).SetMessage("Quotation not found")
	}
	return r.GetQuotation(ctx, tenantID, quotationID)
}

func (r *internalCommerceRepo) GetOrders(ctx context.Context, filter coreentity.InternalCommerceListFilter) ([]coreentity.InternalOrder, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:GetOrders")
	defer span.End()

	filter = normalizeCommerceListFilter(filter)
	type row struct {
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
	where := `WHERE deleted_at IS NULL`
	args := []any{}
	if filter.TenantID != "" {
		where += ` AND tenant_id = ?`
		args = append(args, filter.TenantID)
	}
	var total int
	countQuery := `SELECT COUNT(*) FROM internal_orders ` + where
	if err := r.exec(ctx).GetContext(ctx, &total, r.exec(ctx).Rebind(countQuery), args...); err != nil {
		return nil, 0, err
	}
	rows := make([]row, 0)
	query := `
		SELECT id, tenant_id, company_id, order_number, status, currency_code,
			subtotal_amount, discount_amount, tax_amount, total_amount,
			internal_product_id, internal_product_pricing_id, pricing_snapshot,
			source_type, source_reference_id, metadata, ordered_at, created_at, updated_at
		FROM internal_orders
		` + where + `
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)
	if err := r.exec(ctx).SelectContext(ctx, &rows, r.exec(ctx).Rebind(query), args...); err != nil {
		return nil, 0, err
	}
	items := make([]coreentity.InternalOrder, 0, len(rows))
	for _, item := range rows {
		items = append(items, *mapOrder(item.ID, item.TenantID, item.CompanyID, item.OrderNumber, item.Status, item.CurrencyCode, item.SubtotalAmount, item.DiscountAmount, item.TaxAmount, item.TotalAmount, item.InternalProductID, item.InternalProductPricingID, item.PricingSnapshot, item.SourceType, item.SourceReferenceID, item.Metadata, item.OrderedAt, item.CreatedAt, item.UpdatedAt))
	}
	return items, total, nil
}

func (r *internalCommerceRepo) CreateOrder(ctx context.Context, data coreentity.InternalOrder) (*coreentity.InternalOrder, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:CreateOrder")
	defer span.End()

	number, err := r.nextNumber(ctx, "ORD")
	if err != nil {
		return nil, err
	}
	pricingSnapshot, _ := json.Marshal(data.PricingSnapshot)
	metadata, _ := json.Marshal(data.Metadata)
	query := `
		INSERT INTO internal_orders (
			tenant_id, company_id, order_number, status, currency_code,
			subtotal_amount, discount_amount, tax_amount, total_amount,
			internal_product_id, internal_product_pricing_id, pricing_snapshot,
			source_type, source_reference_id, metadata
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, order_number, ordered_at, created_at
	`
	if err := r.exec(ctx).QueryRowxContext(ctx, r.exec(ctx).Rebind(query),
		data.TenantID, data.CompanyID, number, data.Status, data.CurrencyCode,
		data.SubtotalAmount, data.DiscountAmount, data.TaxAmount, data.TotalAmount,
		data.InternalProductID, data.InternalProductPricingID, pricingSnapshot,
		data.SourceType, data.SourceReferenceID, metadata,
	).Scan(&data.ID, &data.OrderNumber, &data.OrderedAt, &data.CreatedAt); err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *internalCommerceRepo) GetOrder(ctx context.Context, tenantID, id string) (*coreentity.InternalOrder, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:GetOrder")
	defer span.End()

	return r.getOrderByWhere(ctx, tenantID, "o.id = ?", id)
}

func (r *internalCommerceRepo) GetOrderByInvoiceID(ctx context.Context, tenantID, invoiceID string) (*coreentity.InternalOrder, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:GetOrderByInvoiceID")
	defer span.End()

	return r.getOrderByWhere(ctx, tenantID, "i.id = ?", invoiceID)
}

func (r *internalCommerceRepo) getOrderByWhere(ctx context.Context, tenantID, where, value string) (*coreentity.InternalOrder, error) {
	type row struct {
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
	var item row
	query := fmt.Sprintf(`
		SELECT o.id, o.tenant_id, o.company_id, o.order_number, o.status, o.currency_code,
			o.subtotal_amount, o.discount_amount, o.tax_amount, o.total_amount,
			o.internal_product_id, o.internal_product_pricing_id, o.pricing_snapshot,
			o.source_type, o.source_reference_id, o.metadata, o.ordered_at, o.created_at, o.updated_at
		FROM internal_orders o
		LEFT JOIN internal_invoices i ON i.internal_order_id = o.id AND i.deleted_at IS NULL
		WHERE %s AND o.tenant_id = ? AND o.deleted_at IS NULL
	`, where)
	if err := r.exec(ctx).GetContext(ctx, &item, r.exec(ctx).Rebind(query), value, tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Order not found")
		}
		return nil, err
	}
	return mapOrder(item.ID, item.TenantID, item.CompanyID, item.OrderNumber, item.Status, item.CurrencyCode, item.SubtotalAmount, item.DiscountAmount, item.TaxAmount, item.TotalAmount, item.InternalProductID, item.InternalProductPricingID, item.PricingSnapshot, item.SourceType, item.SourceReferenceID, item.Metadata, item.OrderedAt, item.CreatedAt, item.UpdatedAt), nil
}

func (r *internalCommerceRepo) MarkOrderPaid(ctx context.Context, tenantID, id string) (*coreentity.InternalOrder, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:MarkOrderPaid")
	defer span.End()

	query := `UPDATE internal_orders SET status = 'paid', updated_at = NOW() WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`
	if _, err := r.exec(ctx).ExecContext(ctx, r.exec(ctx).Rebind(query), id, tenantID); err != nil {
		return nil, err
	}
	return r.GetOrder(ctx, tenantID, id)
}

func (r *internalCommerceRepo) CreateInvoice(ctx context.Context, data coreentity.InternalInvoice) (*coreentity.InternalInvoice, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:CreateInvoice")
	defer span.End()

	number, err := r.nextNumber(ctx, "INV")
	if err != nil {
		return nil, err
	}
	metadata, _ := json.Marshal(data.Metadata)
	query := `
		INSERT INTO internal_invoices (internal_order_id, invoice_number, status, currency_code, amount, amount_paid, amount_outstanding, due_at, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, invoice_number, created_at
	`
	if err := r.exec(ctx).QueryRowxContext(ctx, r.exec(ctx).Rebind(query), data.InternalOrderID, number, data.Status, data.CurrencyCode, data.Amount, data.AmountPaid, data.AmountOutstanding, data.DueAt, metadata).Scan(&data.ID, &data.InvoiceNumber, &data.CreatedAt); err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *internalCommerceRepo) GetInvoice(ctx context.Context, tenantID, id string) (*coreentity.InternalInvoice, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:GetInvoice")
	defer span.End()

	return r.getInvoiceByWhere(ctx, tenantID, "i.id = ?", id)
}

func (r *internalCommerceRepo) GetInvoiceByOrderID(ctx context.Context, tenantID, orderID string) (*coreentity.InternalInvoice, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:GetInvoiceByOrderID")
	defer span.End()

	return r.getInvoiceByWhere(ctx, tenantID, "i.internal_order_id = ?", orderID)
}

func (r *internalCommerceRepo) GetInvoices(ctx context.Context, filter coreentity.InternalCommerceListFilter) ([]coreentity.InternalInvoice, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:GetInvoices")
	defer span.End()

	filter = normalizeCommerceListFilter(filter)
	type row struct {
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
	where := `WHERE i.deleted_at IS NULL AND o.deleted_at IS NULL`
	args := []any{}
	if filter.TenantID != "" {
		where += ` AND o.tenant_id = ?`
		args = append(args, filter.TenantID)
	}
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM internal_invoices i
		JOIN internal_orders o ON o.id = i.internal_order_id
		` + where
	if err := r.exec(ctx).GetContext(ctx, &total, r.exec(ctx).Rebind(countQuery), args...); err != nil {
		return nil, 0, err
	}
	rows := make([]row, 0)
	query := `
		SELECT i.id, i.internal_order_id, i.invoice_number, i.status, i.currency_code,
			i.amount, i.amount_paid, i.amount_outstanding, i.due_at, i.paid_at, i.expired_at,
			i.metadata, i.created_at, i.updated_at
		FROM internal_invoices i
		JOIN internal_orders o ON o.id = i.internal_order_id
		` + where + `
		ORDER BY i.created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)
	if err := r.exec(ctx).SelectContext(ctx, &rows, r.exec(ctx).Rebind(query), args...); err != nil {
		return nil, 0, err
	}
	items := make([]coreentity.InternalInvoice, 0, len(rows))
	for _, item := range rows {
		items = append(items, *mapInvoice(item.ID, item.InternalOrderID, item.InvoiceNumber, item.Status, item.CurrencyCode, item.Amount, item.AmountPaid, item.AmountOutstanding, item.DueAt, item.PaidAt, item.ExpiredAt, item.Metadata, item.CreatedAt, item.UpdatedAt))
	}
	return items, total, nil
}

func (r *internalCommerceRepo) getInvoiceByWhere(ctx context.Context, tenantID, where, value string) (*coreentity.InternalInvoice, error) {
	type row struct {
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
	var item row
	query := fmt.Sprintf(`
		SELECT i.id, i.internal_order_id, i.invoice_number, i.status, i.currency_code,
			i.amount, i.amount_paid, i.amount_outstanding, i.due_at, i.paid_at, i.expired_at, i.metadata, i.created_at, i.updated_at
		FROM internal_invoices i
		JOIN internal_orders o ON o.id = i.internal_order_id AND o.deleted_at IS NULL
		WHERE %s AND o.tenant_id = ? AND i.deleted_at IS NULL
	`, where)
	if err := r.exec(ctx).GetContext(ctx, &item, r.exec(ctx).Rebind(query), value, tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Invoice not found")
		}
		return nil, err
	}
	return mapInvoice(item.ID, item.InternalOrderID, item.InvoiceNumber, item.Status, item.CurrencyCode, item.Amount, item.AmountPaid, item.AmountOutstanding, item.DueAt, item.PaidAt, item.ExpiredAt, item.Metadata, item.CreatedAt, item.UpdatedAt), nil
}

func (r *internalCommerceRepo) MarkInvoicePaid(ctx context.Context, tenantID, id string, amountPaid string) (*coreentity.InternalInvoice, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:MarkInvoicePaid")
	defer span.End()

	query := `
		UPDATE internal_invoices i
		SET status = 'paid', amount_paid = ?, amount_outstanding = 0, paid_at = NOW(), updated_at = NOW()
		FROM internal_orders o
		WHERE i.internal_order_id = o.id AND i.id = ? AND o.tenant_id = ? AND i.deleted_at IS NULL
	`
	if _, err := r.exec(ctx).ExecContext(ctx, r.exec(ctx).Rebind(query), amountPaid, id, tenantID); err != nil {
		return nil, err
	}
	return r.GetInvoice(ctx, tenantID, id)
}

func (r *internalCommerceRepo) CreatePaymentAttempt(ctx context.Context, data coreentity.InternalPaymentAttempt) (*coreentity.InternalPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:CreatePaymentAttempt")
	defer span.End()

	payload, _ := json.Marshal(data.ProviderPayloadSnapshot)
	metadata, _ := json.Marshal(data.Metadata)
	query := `
		INSERT INTO internal_payment_attempts (
			internal_invoice_id, provider, payment_method_type, payment_channel_code, provider_reference,
			provider_request_id, provider_payment_url, provider_payload_snapshot, status, requested_amount,
			paid_amount, expired_at, raw_last_status, metadata
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at
	`
	if err := r.exec(ctx).QueryRowxContext(ctx, r.exec(ctx).Rebind(query), data.InternalInvoiceID, data.Provider, data.PaymentMethodType, data.PaymentChannelCode, data.ProviderReference, data.ProviderRequestID, data.ProviderPaymentURL, payload, data.Status, data.RequestedAmount, data.PaidAmount, data.ExpiredAt, data.RawLastStatus, metadata).Scan(&data.ID, &data.CreatedAt); err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *internalCommerceRepo) UpdatePaymentAttemptGateway(ctx context.Context, data coreentity.InternalPaymentAttempt) (*coreentity.InternalPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:UpdatePaymentAttemptGateway")
	defer span.End()

	payload, _ := json.Marshal(data.ProviderPayloadSnapshot)
	query := `
		UPDATE internal_payment_attempts
		SET provider_request_id = ?, provider_payment_url = ?, provider_payload_snapshot = ?, status = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
	if _, err := r.exec(ctx).ExecContext(ctx, r.exec(ctx).Rebind(query), data.ProviderRequestID, data.ProviderPaymentURL, payload, data.Status, data.ID); err != nil {
		return nil, err
	}
	return r.GetPaymentAttempt(ctx, "", data.ID)
}

func (r *internalCommerceRepo) GetPaymentAttempt(ctx context.Context, tenantID, id string) (*coreentity.InternalPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:GetPaymentAttempt")
	defer span.End()

	items, err := r.getPaymentAttemptsByWhere(ctx, tenantID, "pa.id = ?", id, 1)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errmsg.NewCustomErrors(404).SetMessage("Payment attempt not found")
	}
	return &items[0], nil
}

func (r *internalCommerceRepo) GetLatestPaymentAttempt(ctx context.Context, tenantID, invoiceID string) (*coreentity.InternalPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:GetLatestPaymentAttempt")
	defer span.End()

	items, err := r.getPaymentAttemptsByWhere(ctx, tenantID, "pa.internal_invoice_id = ?", invoiceID, 1)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &items[0], nil
}

func (r *internalCommerceRepo) GetPaymentAttempts(ctx context.Context, tenantID, invoiceID string) ([]coreentity.InternalPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:GetPaymentAttempts")
	defer span.End()

	return r.getPaymentAttemptsByWhere(ctx, tenantID, "pa.internal_invoice_id = ?", invoiceID, 100)
}

func (r *internalCommerceRepo) getPaymentAttemptsByWhere(ctx context.Context, tenantID, where, value string, limit int) ([]coreentity.InternalPaymentAttempt, error) {
	type row struct {
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
	rows := make([]row, 0)
	args := []any{value}
	query := fmt.Sprintf(`
		SELECT pa.id, pa.internal_invoice_id, pa.provider, pa.payment_method_type, pa.payment_channel_code,
			pa.provider_reference, pa.provider_request_id, pa.provider_payment_url, pa.provider_payload_snapshot,
			pa.status, pa.requested_amount, pa.paid_amount, pa.expired_at, pa.paid_at, pa.failed_at,
			pa.raw_last_status, pa.metadata, pa.created_at, pa.updated_at
		FROM internal_payment_attempts pa
		JOIN internal_invoices i ON i.id = pa.internal_invoice_id AND i.deleted_at IS NULL
		JOIN internal_orders o ON o.id = i.internal_order_id AND o.deleted_at IS NULL
		WHERE %s AND pa.deleted_at IS NULL
	`, where)
	if tenantID != "" {
		query += ` AND o.tenant_id = ?`
		args = append(args, tenantID)
	}
	query += ` ORDER BY pa.created_at DESC LIMIT ?`
	args = append(args, limit)
	if err := r.exec(ctx).SelectContext(ctx, &rows, r.exec(ctx).Rebind(query), args...); err != nil {
		return nil, err
	}
	items := make([]coreentity.InternalPaymentAttempt, 0, len(rows))
	for _, row := range rows {
		items = append(items, *mapPaymentAttempt(row.ID, row.InternalInvoiceID, row.Provider, row.PaymentMethodType, row.PaymentChannelCode, row.ProviderReference, row.ProviderRequestID, row.ProviderPaymentURL, row.ProviderPayloadSnapshot, row.Status, row.RequestedAmount, row.PaidAmount, row.ExpiredAt, row.PaidAt, row.FailedAt, row.RawLastStatus, row.Metadata, row.CreatedAt, row.UpdatedAt))
	}
	return items, nil
}

func (r *internalCommerceRepo) CreatePaymentReceipt(ctx context.Context, data coreentity.InternalPaymentReceipt) (*coreentity.InternalPaymentReceipt, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:CreatePaymentReceipt")
	defer span.End()

	number, err := r.nextNumber(ctx, "RCP")
	if err != nil {
		return nil, err
	}
	metadata, _ := json.Marshal(data.Metadata)
	query := `
		INSERT INTO internal_payment_receipts (
			internal_invoice_id, internal_payment_attempt_id, receipt_number, status, amount_received,
			currency_code, received_at, verified_at, verified_by_user_id, source_type, reference_number, notes, metadata
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, receipt_number, created_at
	`
	if err := r.exec(ctx).QueryRowxContext(ctx, r.exec(ctx).Rebind(query), data.InternalInvoiceID, data.InternalPaymentAttemptID, number, data.Status, data.AmountReceived, data.CurrencyCode, data.ReceivedAt, data.VerifiedAt, data.VerifiedByUserID, data.SourceType, data.ReferenceNumber, data.Notes, metadata).Scan(&data.ID, &data.ReceiptNumber, &data.CreatedAt); err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *internalCommerceRepo) GetPricingSnapshot(ctx context.Context, productID, pricingID, currencyCode string) (map[string]any, string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalcommerce:repo:GetPricingSnapshot")
	defer span.End()

	type row struct {
		ProductID    string          `db:"product_id"`
		ProductCode  string          `db:"product_code"`
		ProductName  string          `db:"product_name"`
		PricingID    string          `db:"pricing_id"`
		PricingCode  string          `db:"pricing_code"`
		PricingName  string          `db:"pricing_name"`
		CurrencyCode string          `db:"currency_code"`
		Amount       decimal.Decimal `db:"amount"`
	}
	var item row
	query := `
		SELECT p.id product_id, p.code product_code, p.name product_name,
			pp.id pricing_id, pp.code pricing_code, pp.name pricing_name,
			price.currency_code, price.amount
		FROM internal_products p
		JOIN internal_product_pricings pp ON pp.internal_product_id = p.id AND pp.deleted_at IS NULL
		JOIN internal_product_prices price ON price.internal_product_pricing_id = pp.id AND price.deleted_at IS NULL
		WHERE p.id = ? AND pp.id = ? AND price.currency_code = ?
			AND p.deleted_at IS NULL AND p.status = 'active' AND pp.status = 'active'
			AND price.started_at <= NOW() AND (price.ended_at IS NULL OR price.ended_at > NOW())
		ORDER BY price.started_at DESC
		LIMIT 1
	`
	if err := r.exec(ctx).GetContext(ctx, &item, r.exec(ctx).Rebind(query), productID, pricingID, currencyCode); err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errmsg.NewCustomErrors(404).SetMessage("Active product pricing not found")
		}
		return nil, "", err
	}
	return map[string]any{
		"product": map[string]any{"id": item.ProductID, "code": item.ProductCode, "name": item.ProductName},
		"pricing": map[string]any{"id": item.PricingID, "code": item.PricingCode, "name": item.PricingName},
		"price":   map[string]any{"currency_code": item.CurrencyCode, "amount": item.Amount.StringFixed(2)},
	}, item.Amount.String(), nil
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

func mapQuotation(id, tenantID string, companyID *string, number, status, currency string, subtotal, discount, tax, total decimal.Decimal, productID, pricingID string, pricingSnapshot, quoteSnapshot json.RawMessage, expiresAt, approvedAt, convertedToOrderID *string, metadata json.RawMessage, createdAt string, updatedAt *string) *coreentity.InternalQuotation {
	item := &coreentity.InternalQuotation{ID: id, TenantID: tenantID, CompanyID: companyID, QuotationNumber: number, Status: status, CurrencyCode: currency, SubtotalAmount: subtotal, DiscountAmount: discount, TaxAmount: tax, TotalAmount: total, InternalProductID: productID, InternalProductPricingID: pricingID, PricingSnapshot: mapJSON(pricingSnapshot), QuoteSnapshot: mapJSON(quoteSnapshot), ExpiresAt: expiresAt, ApprovedAt: approvedAt, ConvertedToOrderID: convertedToOrderID, Metadata: mapJSON(metadata), CreatedAt: createdAt}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}
	return item
}

func mapOrder(id, tenantID string, companyID *string, number, status, currency string, subtotal, discount, tax, total decimal.Decimal, productID, pricingID string, pricingSnapshot json.RawMessage, sourceType string, sourceReferenceID *string, metadata json.RawMessage, orderedAt, createdAt string, updatedAt *string) *coreentity.InternalOrder {
	item := &coreentity.InternalOrder{ID: id, TenantID: tenantID, CompanyID: companyID, OrderNumber: number, Status: status, CurrencyCode: currency, SubtotalAmount: subtotal, DiscountAmount: discount, TaxAmount: tax, TotalAmount: total, InternalProductID: productID, InternalProductPricingID: pricingID, PricingSnapshot: mapJSON(pricingSnapshot), SourceType: sourceType, SourceReferenceID: sourceReferenceID, Metadata: mapJSON(metadata), OrderedAt: orderedAt, CreatedAt: createdAt}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}
	return item
}

func mapInvoice(id, orderID, number, status, currency string, amount, paid, outstanding decimal.Decimal, dueAt, paidAt, expiredAt *string, metadata json.RawMessage, createdAt string, updatedAt *string) *coreentity.InternalInvoice {
	item := &coreentity.InternalInvoice{ID: id, InternalOrderID: orderID, InvoiceNumber: number, Status: status, CurrencyCode: currency, Amount: amount, AmountPaid: paid, AmountOutstanding: outstanding, DueAt: dueAt, PaidAt: paidAt, ExpiredAt: expiredAt, Metadata: mapJSON(metadata), CreatedAt: createdAt}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}
	return item
}

func mapPaymentAttempt(id, invoiceID, provider, methodType, channel, providerReference, requestID, paymentURL string, payload json.RawMessage, status string, requested, paid decimal.Decimal, expiredAt, paidAt, failedAt *string, rawStatus string, metadata json.RawMessage, createdAt string, updatedAt *string) *coreentity.InternalPaymentAttempt {
	item := &coreentity.InternalPaymentAttempt{ID: id, InternalInvoiceID: invoiceID, Provider: provider, PaymentMethodType: methodType, PaymentChannelCode: channel, ProviderReference: providerReference, ProviderRequestID: requestID, ProviderPaymentURL: paymentURL, ProviderPayloadSnapshot: mapJSON(payload), Status: status, RequestedAmount: requested, PaidAmount: paid, ExpiredAt: expiredAt, PaidAt: paidAt, FailedAt: failedAt, RawLastStatus: rawStatus, Metadata: mapJSON(metadata), CreatedAt: createdAt}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}
	return item
}
