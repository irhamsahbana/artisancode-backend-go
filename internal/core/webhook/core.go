package core

import (
	corePorts "codebase-app/internal/ports/core"
	integrationPorts "codebase-app/internal/ports/integration"
)

var _ corePorts.WebhookCore = &webhookCore{}

type webhookCore struct {
	dokuVerifier integrationPorts.DOKUWebhookVerifier
}

type WebhookCoreConfig struct {
	DOKUVerifier integrationPorts.DOKUWebhookVerifier
}

func NewWebhookCore(cfg WebhookCoreConfig) *webhookCore {
	return &webhookCore{
		dokuVerifier: cfg.DOKUVerifier,
	}
}
