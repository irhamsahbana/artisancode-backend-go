package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
)

type internalQuotationCore struct {
	quotationRepo portsRepo.InternalQuotationRepository
	orderRepo     portsRepo.InternalOrderRepository
	invoiceRepo   portsRepo.InternalInvoiceRepository
	tx            portsRepo.Transactor
}

type Config struct {
	QuotationRepo portsRepo.InternalQuotationRepository
	OrderRepo     portsRepo.InternalOrderRepository
	InvoiceRepo   portsRepo.InternalInvoiceRepository
	Tx            portsRepo.Transactor
}

var _ corePorts.InternalQuotationCore = &internalQuotationCore{}

func NewInternalQuotationCore(cfg Config) corePorts.InternalQuotationCore {
	return &internalQuotationCore{
		quotationRepo: cfg.QuotationRepo,
		orderRepo:     cfg.OrderRepo,
		invoiceRepo:   cfg.InvoiceRepo,
		tx:            cfg.Tx,
	}
}
