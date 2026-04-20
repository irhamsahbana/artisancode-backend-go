package consumer

import (
	"context"

	"codebase-app/internal/adapter"
	"codebase-app/internal/setup"
)

func (a *App) build(ctx context.Context) error {
	adapter.Adapters.Sync(
		adapter.WithPostgres(),
		adapter.WithMessagePublisher(),
		adapter.WithMessageSubscriptionManager(),
		adapter.WithEmailSender(),
		adapter.WithStorage(),
	)

	deps, err := setup.NewConsumerDependencies(
		ctx,
		adapter.Adapters.Unsync,
	)
	if err != nil {
		return err
	}

	a.emailSubscription = deps.EmailSubscription
	a.exportSubscription = deps.ExportSubscription
	a.exportCore = deps.ExportCore
	a.exportPublisher = deps.ExportPublisher
	a.subscriptionManager = deps.SubscriptionManager

	if a.shutdown == nil {
		a.shutdown = deps.Shutdown
	}

	return nil
}
