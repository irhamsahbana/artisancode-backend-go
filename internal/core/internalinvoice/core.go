package core

import (
	corePorts "codebase-app/internal/ports/core"
	dokuPorts "codebase-app/internal/ports/integration"
	portsRepo "codebase-app/internal/ports/secondary/db"
)

type internalInvoiceCore struct {
	invoiceRepo   portsRepo.InternalInvoiceRepository
	orderRepo     portsRepo.InternalOrderRepository
	quotationRepo portsRepo.InternalQuotationRepository
	tx            portsRepo.Transactor
	doku          dokuPorts.DokuClient
}

type Config struct {
	InvoiceRepo   portsRepo.InternalInvoiceRepository
	OrderRepo     portsRepo.InternalOrderRepository
	QuotationRepo portsRepo.InternalQuotationRepository
	Tx            portsRepo.Transactor
	DOKU          dokuPorts.DokuClient
}

var _ corePorts.InternalInvoiceCore = &internalInvoiceCore{}

func NewInternalInvoiceCore(cfg Config) corePorts.InternalInvoiceCore {
	return &internalInvoiceCore{
		invoiceRepo:   cfg.InvoiceRepo,
		orderRepo:     cfg.OrderRepo,
		quotationRepo: cfg.QuotationRepo,
		tx:            cfg.Tx,
		doku:          cfg.DOKU,
	}
}
