package core

import (
	corePorts "codebase-app/internal/ports/core"
	integrationPorts "codebase-app/internal/ports/integration"
	repositoryPorts "codebase-app/internal/ports/secondary/db"
)

var _ corePorts.WebhookCore = &webhookCore{}

type webhookCore struct {
	dokuVerifier integrationPorts.DOKUWebhookVerifier
	billingRepo  repositoryPorts.InternalTenantBillingRepository
}

type Config struct {
	DOKUVerifier integrationPorts.DOKUWebhookVerifier
	BillingRepo  repositoryPorts.InternalTenantBillingRepository
}

func NewWebhookCore(cfg Config) *webhookCore {
	return &webhookCore{
		dokuVerifier: cfg.DOKUVerifier,
		billingRepo:  cfg.BillingRepo,
	}
}
