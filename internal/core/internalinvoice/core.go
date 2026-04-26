package core

import (
	"context"
	"strings"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	dokuPorts "codebase-app/internal/ports/integration"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/pkg/errmsg"

	"github.com/shopspring/decimal"
)

type internalInvoiceCore struct {
	repo portsRepo.InternalCommerceRepository
	tx   portsRepo.Transactor
	doku dokuPorts.DokuClient
}

type Config struct {
	Repo portsRepo.InternalCommerceRepository
	Tx   portsRepo.Transactor
	DOKU dokuPorts.DokuClient
}

var _ corePorts.InternalInvoiceCore = &internalInvoiceCore{}

func NewInternalInvoiceCore(cfg Config) corePorts.InternalInvoiceCore {
	return &internalInvoiceCore{repo: cfg.Repo, tx: cfg.Tx, doku: cfg.DOKU}
}

func (c *internalInvoiceCore) GetInvoices(ctx context.Context, filter coreentity.InternalCommerceListFilter) ([]coreentity.InternalInvoice, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalinvoice:core:GetInvoices")
	defer span.End()

	if filter.TenantID == "" {
		filter.TenantID = common.GetUserContext(ctx).TenantID
	}
	return c.repo.GetInvoices(ctx, filter)
}

func (c *internalInvoiceCore) GetInvoice(ctx context.Context, id string) (*coreentity.InternalInvoiceDetail, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalinvoice:core:GetInvoice")
	defer span.End()

	tenantID := common.GetUserContext(ctx).TenantID
	invoice, err := c.repo.GetInvoice(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	order, err := c.repo.GetOrderByInvoiceID(ctx, tenantID, invoice.ID)
	if err != nil {
		return nil, err
	}
	var quote *coreentity.InternalQuotation
	if order.SourceType == coreentity.OrderSourceTypeQuotation && order.SourceReferenceID != nil {
		quote, _ = c.repo.GetQuotation(ctx, tenantID, *order.SourceReferenceID)
	}
	latest, _ := c.repo.GetLatestPaymentAttempt(ctx, tenantID, invoice.ID)
	return &coreentity.InternalInvoiceDetail{Invoice: *invoice, Order: *order, Quotation: quote, LatestPaymentAttempt: latest}, nil
}

func (c *internalInvoiceCore) ExecuteInvoiceAction(ctx context.Context, input coreentity.InternalInvoiceActionInput) (*coreentity.InternalPaymentAttemptResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalinvoice:core:ExecuteInvoiceAction")
	defer span.End()

	switch strings.TrimSpace(input.Action) {
	case coreentity.ActionCreateDOKUPaymentAttempt:
		return c.createDOKUAttempt(ctx, input)
	case coreentity.ActionCreateCustomPaymentAttempt:
		return c.createManualAttempt(ctx, input)
	default:
		return nil, errmsg.NewCustomErrors(400).SetMessage("Unsupported invoice action")
	}
}

func (c *internalInvoiceCore) GetPaymentAttempts(ctx context.Context, invoiceID string) ([]coreentity.InternalPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalinvoice:core:GetPaymentAttempts")
	defer span.End()

	return c.repo.GetPaymentAttempts(ctx, common.GetUserContext(ctx).TenantID, invoiceID)
}

func (c *internalInvoiceCore) ExecutePaymentAttemptAction(ctx context.Context, input coreentity.InternalPaymentAttemptActionInput) (*coreentity.InternalPaymentAttemptResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalinvoice:core:ExecutePaymentAttemptAction")
	defer span.End()

	if input.Action != coreentity.ActionRetryPayment {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Unsupported payment attempt action")
	}
	attempt, err := c.repo.GetPaymentAttempt(ctx, input.UserCtx.TenantID, input.ID)
	if err != nil {
		return nil, err
	}
	if attempt.Status != coreentity.PaymentAttemptStatusFailed && attempt.Status != coreentity.PaymentAttemptStatusExpired && attempt.Status != coreentity.PaymentAttemptStatusCancelled {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Payment attempt cannot be retried")
	}
	return c.ExecuteInvoiceAction(ctx, coreentity.InternalInvoiceActionInput{
		UserCtx: input.UserCtx, ID: attempt.InternalInvoiceID, Action: actionForProvider(input.Provider, attempt.Provider),
		Provider: input.Provider, PaymentMethodType: input.PaymentMethodType, PaymentChannelCode: input.PaymentChannelCode,
		CallbackURL: input.CallbackURL, CustomerName: input.CustomerName, CustomerEmail: input.CustomerEmail, CustomerPhone: input.CustomerPhone,
		RequestID: input.RequestID,
	})
}

func (c *internalInvoiceCore) CreatePaymentReceipt(ctx context.Context, receipt coreentity.InternalPaymentReceipt) (*coreentity.InternalPaymentReceiptResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalinvoice:core:CreatePaymentReceipt")
	defer span.End()

	tenantID := receipt.UserCtx.TenantID
	if receipt.Status == "" {
		receipt.Status = coreentity.PaymentReceiptStatusPendingVerification
	}
	receipt.SourceType = coreentity.PaymentProviderManual
	if receipt.Status == coreentity.PaymentReceiptStatusAccepted {
		now := time.Now().UTC().Format(time.RFC3339)
		receipt.VerifiedAt = &now
		receipt.VerifiedByUserID = &receipt.UserCtx.UserID
	}
	var result coreentity.InternalPaymentReceiptResult
	err := c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		invoice, err := c.repo.GetInvoice(txCtx, tenantID, receipt.InternalInvoiceID)
		if err != nil {
			return err
		}
		if receipt.AmountReceived.LessThan(invoice.AmountOutstanding) && receipt.Status == coreentity.PaymentReceiptStatusAccepted {
			return errmsg.NewCustomErrors(400).SetMessage("Payment amount does not match invoice outstanding amount")
		}
		created, err := c.repo.CreatePaymentReceipt(txCtx, receipt)
		if err != nil {
			return err
		}
		order, err := c.repo.GetOrderByInvoiceID(txCtx, tenantID, invoice.ID)
		if err != nil {
			return err
		}
		if receipt.Status == coreentity.PaymentReceiptStatusAccepted {
			invoice, err = c.repo.MarkInvoicePaid(txCtx, tenantID, invoice.ID, receipt.AmountReceived.String())
			if err != nil {
				return err
			}
			order, err = c.repo.MarkOrderPaid(txCtx, tenantID, order.ID)
			if err != nil {
				return err
			}
		}
		result = coreentity.InternalPaymentReceiptResult{Receipt: *created, Invoice: *invoice, Order: *order}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *internalInvoiceCore) createDOKUAttempt(ctx context.Context, input coreentity.InternalInvoiceActionInput) (*coreentity.InternalPaymentAttemptResult, error) {
	if c.doku == nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage("DOKU client is not configured")
	}
	invoice, err := c.repo.GetInvoice(ctx, input.UserCtx.TenantID, input.ID)
	if err != nil {
		return nil, err
	}
	if invoice.Status != coreentity.InvoiceStatusOpen && invoice.Status != coreentity.InvoiceStatusPartiallyPaid {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Invoice is not payable")
	}
	attempt, err := c.repo.CreatePaymentAttempt(ctx, coreentity.InternalPaymentAttempt{
		InternalInvoiceID: invoice.ID, Provider: coreentity.PaymentProviderDOKU, PaymentMethodType: "gateway",
		ProviderReference: invoice.InvoiceNumber, Status: coreentity.PaymentAttemptStatusInitiated,
		RequestedAmount: invoice.AmountOutstanding, PaidAmount: decimal.Zero, Metadata: map[string]any{},
	})
	if err != nil {
		return nil, err
	}
	resp, err := c.doku.CreatePayment(ctx, restentity.DokuCreatePaymentRequest{
		InvoiceNumber: invoice.InvoiceNumber, Amount: invoice.AmountOutstanding.Round(0).IntPart(),
		CustomerName: input.CustomerName, CustomerEmail: input.CustomerEmail, CustomerPhone: input.CustomerPhone,
		CallbackURL: input.CallbackURL,
	})
	if err != nil {
		return nil, err
	}
	attempt.ProviderRequestID = resp.RequestID
	attempt.ProviderPaymentURL = resp.PaymentURL
	attempt.ProviderPayloadSnapshot = map[string]any{"doku_response": resp}
	attempt.Status = coreentity.PaymentAttemptStatusPending
	attempt, err = c.repo.UpdatePaymentAttemptGateway(ctx, *attempt)
	if err != nil {
		return nil, err
	}
	return &coreentity.InternalPaymentAttemptResult{PaymentAttempt: *attempt}, nil
}

func (c *internalInvoiceCore) createManualAttempt(ctx context.Context, input coreentity.InternalInvoiceActionInput) (*coreentity.InternalPaymentAttemptResult, error) {
	invoice, err := c.repo.GetInvoice(ctx, input.UserCtx.TenantID, input.ID)
	if err != nil {
		return nil, err
	}
	if invoice.Status != coreentity.InvoiceStatusOpen && invoice.Status != coreentity.InvoiceStatusPartiallyPaid {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Invoice is not payable")
	}
	provider := input.Provider
	if provider == "" {
		provider = coreentity.PaymentProviderManual
	}
	method := input.PaymentMethodType
	if method == "" {
		method = "bank_transfer"
	}
	channel := input.PaymentChannelCode
	if channel == "" {
		channel = "manual_bank_transfer"
	}
	attempt, err := c.repo.CreatePaymentAttempt(ctx, coreentity.InternalPaymentAttempt{
		InternalInvoiceID: invoice.ID, Provider: provider, PaymentMethodType: method, PaymentChannelCode: channel,
		ProviderReference: invoice.InvoiceNumber, Status: coreentity.PaymentAttemptStatusPending,
		RequestedAmount: invoice.AmountOutstanding, PaidAmount: decimal.Zero, Metadata: map[string]any{},
	})
	if err != nil {
		return nil, err
	}
	return &coreentity.InternalPaymentAttemptResult{
		PaymentAttempt: *attempt,
		Instruction: &coreentity.InternalManualPaymentInstruction{
			Title: "Manual transfer", ReferenceNumber: invoice.InvoiceNumber, BankName: "BCA",
			AccountNumber: "1234567890", AccountName: "PT Artisan Code", Notes: "Use invoice number as transfer note",
		},
	}, nil
}

func actionForProvider(requestProvider, previousProvider string) string {
	provider := requestProvider
	if provider == "" {
		provider = previousProvider
	}
	if provider == coreentity.PaymentProviderDOKU {
		return coreentity.ActionCreateDOKUPaymentAttempt
	}
	return coreentity.ActionCreateCustomPaymentAttempt
}
