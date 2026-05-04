package core

import (
	dokuPorts "codebase-app/internal/ports/integration"
	repositoryPorts "codebase-app/internal/ports/secondary/db"
)

type internalTenantBillingCore struct {
	productRepo  repositoryPorts.InternalProductRepository
	currencyRepo repositoryPorts.InternalCurrencyRepository
	billingRepo  repositoryPorts.InternalTenantBillingRepository
	tx           repositoryPorts.Transactor
	doku         dokuPorts.DokuClient
}

type Config struct {
	ProductRepo  repositoryPorts.InternalProductRepository
	CurrencyRepo repositoryPorts.InternalCurrencyRepository
	BillingRepo  repositoryPorts.InternalTenantBillingRepository
	Tx           repositoryPorts.Transactor
	DOKU         dokuPorts.DokuClient
}

func NewInternalTenantBillingCore(cfg Config) *internalTenantBillingCore {
	return &internalTenantBillingCore{
		productRepo:  cfg.ProductRepo,
		currencyRepo: cfg.CurrencyRepo,
		billingRepo:  cfg.BillingRepo,
		tx:           cfg.Tx,
		doku:         cfg.DOKU,
	}
}
