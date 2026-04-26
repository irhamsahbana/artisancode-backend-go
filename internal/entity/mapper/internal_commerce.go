package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"

	"github.com/shopspring/decimal"
)

func CreateInternalQuotationFromRest(ctx context.Context, req restentity.CreateInternalQuotationReq) coreentity.CreateInternalQuotationInput {
	return coreentity.CreateInternalQuotationInput{
		UserCtx: common.GetUserContext(ctx), TenantID: req.TenantID, InternalProductID: req.InternalProductID, InternalProductPricingID: req.InternalProductPricingID,
		CurrencyCode: req.CurrencyCode, SubtotalAmount: mustDecimal(req.SubtotalAmount), DiscountAmount: mustDecimal(req.DiscountAmount),
		TaxAmount: mustDecimal(req.TaxAmount), TotalAmount: mustDecimal(req.TotalAmount), ExpiresAt: req.ExpiresAt,
		QuoteSnapshot: req.QuoteSnapshot, Metadata: req.Metadata,
	}
}

func CreateInternalOrderFromRest(ctx context.Context, req restentity.CreateInternalOrderReq) coreentity.CreateInternalOrderInput {
	return coreentity.CreateInternalOrderInput{
		UserCtx: common.GetUserContext(ctx), TenantID: req.TenantID, InternalProductID: req.InternalProductID, InternalProductPricingID: req.InternalProductPricingID,
		CurrencyCode: req.CurrencyCode, InvoiceDueAt: req.InvoiceDueAt, Metadata: req.Metadata,
	}
}

func InternalQuotationActionFromRest(ctx context.Context, req restentity.ExecuteQuotationActionReq, requestID string) coreentity.InternalQuotationActionInput {
	return coreentity.InternalQuotationActionInput{
		UserCtx: common.GetUserContext(ctx), ID: req.ID, Action: req.Action, ApprovalNote: req.ApprovalNote, InvoiceDueAt: req.InvoiceDueAt, RequestID: requestID,
	}
}

func InternalInvoiceActionFromRest(ctx context.Context, req restentity.ExecuteInvoiceActionReq, requestID string) coreentity.InternalInvoiceActionInput {
	return coreentity.InternalInvoiceActionInput{
		UserCtx: common.GetUserContext(ctx), ID: req.ID, Action: req.Action, Provider: req.Provider, PaymentMethodType: req.PaymentMethodType,
		PaymentChannelCode: req.PaymentChannelCode, CallbackURL: req.CallbackURL, CustomerName: req.Customer.Name,
		CustomerEmail: req.Customer.Email, CustomerPhone: req.Customer.Phone, RequestID: requestID,
	}
}

func InternalPaymentAttemptActionFromRest(ctx context.Context, req restentity.ExecutePaymentAttemptActionReq, requestID string) coreentity.InternalPaymentAttemptActionInput {
	return coreentity.InternalPaymentAttemptActionInput{
		UserCtx: common.GetUserContext(ctx), ID: req.ID, Action: req.Action, Provider: req.Provider, PaymentMethodType: req.PaymentMethodType,
		PaymentChannelCode: req.PaymentChannelCode, CallbackURL: req.CallbackURL, CustomerName: req.Customer.Name,
		CustomerEmail: req.Customer.Email, CustomerPhone: req.Customer.Phone, RequestID: requestID,
	}
}

func InternalPaymentReceiptFromRest(ctx context.Context, req restentity.CreatePaymentReceiptReq) coreentity.InternalPaymentReceipt {
	return coreentity.InternalPaymentReceipt{
		UserCtx: common.GetUserContext(ctx), InternalInvoiceID: req.InternalInvoiceID, InternalPaymentAttemptID: req.InternalPaymentAttemptID,
		Status: req.Status, AmountReceived: mustDecimal(req.AmountReceived), CurrencyCode: req.CurrencyCode,
		ReceivedAt: req.ReceivedAt, ReferenceNumber: req.ReferenceNumber, Notes: req.Notes, Metadata: req.Metadata,
	}
}

func InternalQuotationToRest(item coreentity.InternalQuotation) restentity.InternalQuotation {
	return restentity.InternalQuotation{
		ID: item.ID, QuotationNumber: item.QuotationNumber, Status: item.Status, CurrencyCode: item.CurrencyCode,
		SubtotalAmount: item.SubtotalAmount.StringFixed(2), DiscountAmount: item.DiscountAmount.StringFixed(2),
		TaxAmount: item.TaxAmount.StringFixed(2), TotalAmount: item.TotalAmount.StringFixed(2),
		InternalProductID: item.InternalProductID, InternalProductPricingID: item.InternalProductPricingID,
		PricingSnapshot: item.PricingSnapshot, QuoteSnapshot: item.QuoteSnapshot, ExpiresAt: item.ExpiresAt,
		ApprovedAt: item.ApprovedAt, ConvertedToOrderID: item.ConvertedToOrderID, Metadata: item.Metadata,
		CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt, AvailableActions: QuotationActions(item),
	}
}

func InternalOrderToRest(item coreentity.InternalOrder, invoice *coreentity.InternalInvoice) restentity.InternalOrder {
	resp := restentity.InternalOrder{
		ID: item.ID, OrderNumber: item.OrderNumber, Status: item.Status, CurrencyCode: item.CurrencyCode,
		SubtotalAmount: item.SubtotalAmount.StringFixed(2), DiscountAmount: item.DiscountAmount.StringFixed(2),
		TaxAmount: item.TaxAmount.StringFixed(2), TotalAmount: item.TotalAmount.StringFixed(2),
		InternalProductID: item.InternalProductID, InternalProductPricingID: item.InternalProductPricingID,
		PricingSnapshot: item.PricingSnapshot, SourceType: item.SourceType, SourceReferenceID: item.SourceReferenceID,
		Metadata: item.Metadata, OrderedAt: item.OrderedAt, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
	}
	if invoice != nil {
		resp.Invoice = &restentity.InternalInvoiceSummary{ID: invoice.ID, InvoiceNumber: invoice.InvoiceNumber, Status: invoice.Status}
		resp.AvailableActions = []restentity.InternalCommerceAction{{Key: coreentity.ActionViewInvoice, Label: "View invoice", Method: "GET", Href: "/internal-commerce/invoices/" + invoice.ID}}
	}
	return resp
}

func InternalInvoiceToRest(detail coreentity.InternalInvoiceDetail) restentity.InternalInvoice {
	invoice := detail.Invoice
	resp := restentity.InternalInvoice{
		ID: invoice.ID, InternalOrderID: invoice.InternalOrderID, InvoiceNumber: invoice.InvoiceNumber, Status: invoice.Status,
		CurrencyCode: invoice.CurrencyCode, Amount: invoice.Amount.StringFixed(2), AmountPaid: invoice.AmountPaid.StringFixed(2),
		AmountOutstanding: invoice.AmountOutstanding.StringFixed(2), DueAt: invoice.DueAt, PaidAt: invoice.PaidAt,
		ExpiredAt: invoice.ExpiredAt, Metadata: invoice.Metadata, CreatedAt: invoice.CreatedAt, UpdatedAt: invoice.UpdatedAt,
		AvailablePaymentMethods: []restentity.InternalPaymentMethod{
			{Provider: coreentity.PaymentProviderDOKU, PaymentMethodType: "gateway", Label: "DOKU"},
			{Provider: coreentity.PaymentProviderManual, PaymentMethodType: "bank_transfer", Label: "Manual transfer"},
		},
		AvailableActions: InvoiceActions(invoice),
	}
	order := InternalOrderToRest(detail.Order, nil)
	resp.Order = &order
	if detail.Quotation != nil {
		resp.Quotation = &restentity.InternalQuotationSummary{ID: detail.Quotation.ID, QuotationNumber: detail.Quotation.QuotationNumber, Status: detail.Quotation.Status}
	}
	if detail.LatestPaymentAttempt != nil {
		attempt := InternalPaymentAttemptToRest(*detail.LatestPaymentAttempt, nil)
		resp.LatestPaymentAttempt = &attempt
	}
	return resp
}

func InternalPaymentAttemptToRest(item coreentity.InternalPaymentAttempt, instruction *coreentity.InternalManualPaymentInstruction) restentity.InternalPaymentAttempt {
	resp := restentity.InternalPaymentAttempt{
		ID: item.ID, InternalInvoiceID: item.InternalInvoiceID, Provider: item.Provider, PaymentMethodType: item.PaymentMethodType,
		PaymentChannelCode: item.PaymentChannelCode, ProviderReference: item.ProviderReference, ProviderRequestID: item.ProviderRequestID,
		PaymentURL: item.ProviderPaymentURL, Status: item.Status, RequestedAmount: item.RequestedAmount.StringFixed(2),
		PaidAmount: item.PaidAmount.StringFixed(2), ExpiredAt: item.ExpiredAt, PaidAt: item.PaidAt, FailedAt: item.FailedAt,
		RawLastStatus: item.RawLastStatus, Metadata: item.Metadata, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
		AvailableActions: PaymentAttemptActions(item),
	}
	if instruction != nil {
		resp.Instruction = &restentity.ManualPaymentInstruction{
			Title: instruction.Title, ReferenceNumber: instruction.ReferenceNumber, BankName: instruction.BankName,
			AccountNumber: instruction.AccountNumber, AccountName: instruction.AccountName, Notes: instruction.Notes,
		}
	}
	return resp
}

func InternalCommerceBundleToRest(bundle coreentity.InternalCommerceBundle) restentity.InternalCommerceBundleResp {
	resp := restentity.InternalCommerceBundleResp{}
	if bundle.Quotation != nil {
		q := InternalQuotationToRest(*bundle.Quotation)
		resp.Quotation = &q
	}
	if bundle.Order != nil {
		o := InternalOrderToRest(*bundle.Order, bundle.Invoice)
		resp.Order = &o
	}
	if bundle.Invoice != nil {
		i := restentity.InternalInvoice{
			ID: bundle.Invoice.ID, InternalOrderID: bundle.Invoice.InternalOrderID, InvoiceNumber: bundle.Invoice.InvoiceNumber,
			Status: bundle.Invoice.Status, CurrencyCode: bundle.Invoice.CurrencyCode, Amount: bundle.Invoice.Amount.StringFixed(2),
			AmountPaid: bundle.Invoice.AmountPaid.StringFixed(2), AmountOutstanding: bundle.Invoice.AmountOutstanding.StringFixed(2),
			DueAt: bundle.Invoice.DueAt, PaidAt: bundle.Invoice.PaidAt, ExpiredAt: bundle.Invoice.ExpiredAt, Metadata: bundle.Invoice.Metadata,
			CreatedAt: bundle.Invoice.CreatedAt, UpdatedAt: bundle.Invoice.UpdatedAt, AvailableActions: InvoiceActions(*bundle.Invoice),
		}
		resp.Invoice = &i
		resp.AvailableActions = i.AvailableActions
	}
	return resp
}

func QuotationActions(item coreentity.InternalQuotation) []restentity.InternalCommerceAction {
	if item.Status == coreentity.QuotationStatusDraft || item.Status == coreentity.QuotationStatusSent {
		return []restentity.InternalCommerceAction{{Key: coreentity.ActionApproveQuotation, Label: "Approve quotation", Method: "POST", Href: "/internal-commerce/quotations/" + item.ID + "/actions"}}
	}
	if item.Status == coreentity.QuotationStatusApproved {
		return []restentity.InternalCommerceAction{{Key: coreentity.ActionConvertQuotation, Label: "Convert to order", Method: "POST", Href: "/internal-commerce/quotations/" + item.ID + "/actions"}}
	}
	return nil
}

func InvoiceActions(item coreentity.InternalInvoice) []restentity.InternalCommerceAction {
	if item.Status != coreentity.InvoiceStatusOpen && item.Status != coreentity.InvoiceStatusPartiallyPaid {
		return nil
	}
	return []restentity.InternalCommerceAction{
		{Key: coreentity.ActionCreateDOKUPaymentAttempt, Label: "Pay with DOKU", Method: "POST", Href: "/internal-commerce/invoices/" + item.ID + "/actions"},
		{Key: coreentity.ActionCreateCustomPaymentAttempt, Label: "Manual transfer", Method: "POST", Href: "/internal-commerce/invoices/" + item.ID + "/actions"},
	}
}

func PaymentAttemptActions(item coreentity.InternalPaymentAttempt) []restentity.InternalCommerceAction {
	if item.Status == coreentity.PaymentAttemptStatusPending && item.ProviderPaymentURL != "" {
		return []restentity.InternalCommerceAction{{Key: coreentity.ActionRedirectToPaymentURL, Label: "Continue payment", Method: "GET", Href: item.ProviderPaymentURL}}
	}
	if item.Status == coreentity.PaymentAttemptStatusFailed || item.Status == coreentity.PaymentAttemptStatusExpired || item.Status == coreentity.PaymentAttemptStatusCancelled {
		return []restentity.InternalCommerceAction{{Key: coreentity.ActionRetryPayment, Label: "Retry payment", Method: "POST", Href: "/internal-commerce/payment-attempts/" + item.ID + "/actions"}}
	}
	if item.Provider == coreentity.PaymentProviderManual && item.Status == coreentity.PaymentAttemptStatusPending {
		return []restentity.InternalCommerceAction{{Key: coreentity.ActionWaitForVerification, Label: "Waiting for verification", Method: "GET", Href: "/internal-commerce/invoices/" + item.InternalInvoiceID}}
	}
	return nil
}

func mustDecimal(value string) decimal.Decimal {
	d, _ := decimal.NewFromString(value)
	return d
}
