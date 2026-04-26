package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/pkg/errmsg"

	"github.com/shopspring/decimal"
)

type internalQuotationCore struct {
	repo portsRepo.InternalCommerceRepository
	tx   portsRepo.Transactor
}

type Config struct {
	Repo portsRepo.InternalCommerceRepository
	Tx   portsRepo.Transactor
}

var _ corePorts.InternalQuotationCore = &internalQuotationCore{}

func NewInternalQuotationCore(cfg Config) corePorts.InternalQuotationCore {
	return &internalQuotationCore{repo: cfg.Repo, tx: cfg.Tx}
}

func (c *internalQuotationCore) GetQuotations(ctx context.Context, filter coreentity.InternalCommerceListFilter) ([]coreentity.InternalQuotation, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalquotation:core:GetQuotations")
	defer span.End()

	if filter.TenantID == "" {
		filter.TenantID = userTenant(ctx)
	}
	return c.repo.GetQuotations(ctx, filter)
}

func (c *internalQuotationCore) CreateQuotation(ctx context.Context, input coreentity.CreateInternalQuotationInput) (*coreentity.InternalQuotation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalquotation:core:CreateQuotation")
	defer span.End()

	tenantID := strings.TrimSpace(input.UserCtx.TenantID)
	if tenantID == "" {
		return nil, errmsg.NewCustomErrors(401).SetMessage("User context is required")
	}
	input.CurrencyCode = strings.ToUpper(strings.TrimSpace(input.CurrencyCode))
	if input.TotalAmount.LessThan(decimal.Zero) || input.SubtotalAmount.LessThan(decimal.Zero) || input.DiscountAmount.LessThan(decimal.Zero) || input.TaxAmount.LessThan(decimal.Zero) {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Amounts must not be negative")
	}
	pricingSnapshot, _, err := c.repo.GetPricingSnapshot(ctx, input.InternalProductID, input.InternalProductPricingID, input.CurrencyCode)
	if err != nil {
		return nil, err
	}
	return c.repo.CreateQuotation(ctx, coreentity.InternalQuotation{
		TenantID: input.UserCtx.TenantID, CompanyID: input.UserCtx.CompanyID, Status: coreentity.QuotationStatusDraft,
		CurrencyCode: input.CurrencyCode, SubtotalAmount: input.SubtotalAmount, DiscountAmount: input.DiscountAmount,
		TaxAmount: input.TaxAmount, TotalAmount: input.TotalAmount, InternalProductID: input.InternalProductID,
		InternalProductPricingID: input.InternalProductPricingID, PricingSnapshot: pricingSnapshot, QuoteSnapshot: input.QuoteSnapshot,
		ExpiresAt: input.ExpiresAt, Metadata: input.Metadata,
	})
}

func (c *internalQuotationCore) GetQuotation(ctx context.Context, id string) (*coreentity.InternalQuotation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalquotation:core:GetQuotation")
	defer span.End()

	return c.repo.GetQuotation(ctx, userTenant(ctx), id)
}

func (c *internalQuotationCore) ExecuteQuotationAction(ctx context.Context, input coreentity.InternalQuotationActionInput) (*coreentity.InternalCommerceBundle, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalquotation:core:ExecuteQuotationAction")
	defer span.End()

	input.Action = strings.TrimSpace(input.Action)
	switch input.Action {
	case coreentity.ActionApproveQuotation:
		quote, err := c.repo.GetQuotation(ctx, input.UserCtx.TenantID, input.ID)
		if err != nil {
			return nil, err
		}
		if quote.Status != coreentity.QuotationStatusDraft && quote.Status != coreentity.QuotationStatusSent {
			return nil, errmsg.NewCustomErrors(400).SetMessage("Quotation cannot be approved from current status")
		}
		approved, err := c.repo.ApproveQuotation(ctx, input.UserCtx.TenantID, input.ID)
		if err != nil {
			return nil, err
		}
		return &coreentity.InternalCommerceBundle{Quotation: approved}, nil
	case coreentity.ActionConvertQuotation:
		var bundle coreentity.InternalCommerceBundle
		err := c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
			quote, err := c.repo.GetQuotation(txCtx, input.UserCtx.TenantID, input.ID)
			if err != nil {
				return err
			}
			if quote.Status != coreentity.QuotationStatusApproved {
				return errmsg.NewCustomErrors(400).SetMessage("Only approved quotation can be converted")
			}
			order, err := c.repo.CreateOrder(txCtx, coreentity.InternalOrder{
				TenantID: quote.TenantID, CompanyID: quote.CompanyID, Status: coreentity.OrderStatusPendingPayment, CurrencyCode: quote.CurrencyCode,
				SubtotalAmount: quote.SubtotalAmount, DiscountAmount: quote.DiscountAmount, TaxAmount: quote.TaxAmount, TotalAmount: quote.TotalAmount,
				InternalProductID: quote.InternalProductID, InternalProductPricingID: quote.InternalProductPricingID, PricingSnapshot: quote.PricingSnapshot,
				SourceType: coreentity.OrderSourceTypeQuotation, SourceReferenceID: &quote.ID, Metadata: quote.Metadata,
			})
			if err != nil {
				return err
			}
			invoice, err := c.repo.CreateInvoice(txCtx, coreentity.InternalInvoice{
				InternalOrderID: order.ID, Status: coreentity.InvoiceStatusOpen, CurrencyCode: order.CurrencyCode,
				Amount: order.TotalAmount, AmountPaid: decimal.Zero, AmountOutstanding: order.TotalAmount, DueAt: input.InvoiceDueAt,
				Metadata: map[string]any{},
			})
			if err != nil {
				return err
			}
			converted, err := c.repo.MarkQuotationConverted(txCtx, quote.TenantID, quote.ID, order.ID)
			if err != nil {
				return err
			}
			bundle.Quotation = converted
			bundle.Order = order
			bundle.Invoice = invoice
			return nil
		})
		if err != nil {
			return nil, err
		}
		return &bundle, nil
	default:
		return nil, errmsg.NewCustomErrors(400).SetMessage("Unsupported quotation action")
	}
}

func userTenant(ctx context.Context) string {
	return coreentityTenant(ctx)
}
