package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
)

type internalOrderCore struct {
	orderRepo   portsRepo.InternalOrderRepository
	invoiceRepo portsRepo.InternalInvoiceRepository
	tx          portsRepo.Transactor
}

type Config struct {
	OrderRepo   portsRepo.InternalOrderRepository
	InvoiceRepo portsRepo.InternalInvoiceRepository
	Tx          portsRepo.Transactor
}

var _ corePorts.InternalOrderCore = &internalOrderCore{}

func NewInternalOrderCore(cfg Config) corePorts.InternalOrderCore {
	return &internalOrderCore{
		orderRepo:   cfg.OrderRepo,
		invoiceRepo: cfg.InvoiceRepo,
		tx:          cfg.Tx,
	}
}
