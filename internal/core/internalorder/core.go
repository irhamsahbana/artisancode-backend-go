package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/pkg/errmsg"

	"github.com/shopspring/decimal"
)

type internalOrderCore struct {
	repo portsRepo.InternalCommerceRepository
	tx   portsRepo.Transactor
}

type Config struct {
	Repo portsRepo.InternalCommerceRepository
	Tx   portsRepo.Transactor
}

var _ corePorts.InternalOrderCore = &internalOrderCore{}

func NewInternalOrderCore(cfg Config) corePorts.InternalOrderCore {
	return &internalOrderCore{repo: cfg.Repo, tx: cfg.Tx}
}

func (c *internalOrderCore) GetOrders(ctx context.Context, filter coreentity.InternalCommerceListFilter) ([]coreentity.InternalOrder, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalorder:core:GetOrders")
	defer span.End()

	if filter.TenantID == "" {
		filter.TenantID = common.GetUserContext(ctx).TenantID
	}
	return c.repo.GetOrders(ctx, filter)
}

func (c *internalOrderCore) CreateOrder(ctx context.Context, input coreentity.CreateInternalOrderInput) (*coreentity.InternalCommerceBundle, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalorder:core:CreateOrder")
	defer span.End()

	uc := input.UserCtx
	if strings.TrimSpace(uc.TenantID) == "" {
		return nil, errmsg.NewCustomErrors(401).SetMessage("User context is required")
	}
	currency := strings.ToUpper(strings.TrimSpace(input.CurrencyCode))
	pricingSnapshot, amount, err := c.repo.GetPricingSnapshot(ctx, input.InternalProductID, input.InternalProductPricingID, currency)
	if err != nil {
		return nil, err
	}
	total, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, err
	}
	var bundle coreentity.InternalCommerceBundle
	err = c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		order, err := c.repo.CreateOrder(txCtx, coreentity.InternalOrder{
			TenantID: uc.TenantID, CompanyID: uc.CompanyID, Status: coreentity.OrderStatusPendingPayment, CurrencyCode: currency,
			SubtotalAmount: total, DiscountAmount: decimal.Zero, TaxAmount: decimal.Zero, TotalAmount: total,
			InternalProductID: input.InternalProductID, InternalProductPricingID: input.InternalProductPricingID, PricingSnapshot: pricingSnapshot,
			SourceType: coreentity.OrderSourceTypeStandardPricing, Metadata: input.Metadata,
		})
		if err != nil {
			return err
		}
		invoice, err := c.repo.CreateInvoice(txCtx, coreentity.InternalInvoice{
			InternalOrderID: order.ID, Status: coreentity.InvoiceStatusOpen, CurrencyCode: currency,
			Amount: total, AmountPaid: decimal.Zero, AmountOutstanding: total, DueAt: input.InvoiceDueAt, Metadata: map[string]any{},
		})
		if err != nil {
			return err
		}
		bundle.Order = order
		bundle.Invoice = invoice
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (c *internalOrderCore) GetOrder(ctx context.Context, id string) (*coreentity.InternalCommerceBundle, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalorder:core:GetOrder")
	defer span.End()

	tenantID := common.GetUserContext(ctx).TenantID
	order, err := c.repo.GetOrder(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	invoice, _ := c.repo.GetInvoiceByOrderID(ctx, tenantID, id)
	return &coreentity.InternalCommerceBundle{Order: order, Invoice: invoice}, nil
}
