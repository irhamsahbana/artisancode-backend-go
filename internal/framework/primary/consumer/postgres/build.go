package consumer

import (
	"context"
	"time"

	"codebase-app/internal/adapter"
	exportjobconsumer "codebase-app/internal/framework/primary/consumer/postgres/export_job"
	userconsumer "codebase-app/internal/framework/primary/consumer/postgres/user"
	userinvitationconsumer "codebase-app/internal/framework/primary/consumer/postgres/userinvitation"
	postgresBus "codebase-app/internal/framework/secondary/publisher/postgres"
	"codebase-app/internal/infrastructure/config"
	"codebase-app/internal/setup"

	"github.com/rs/zerolog/log"
)

func (a *App) build(ctx context.Context) error {
	adapter.Adapters.Sync(
		adapter.WithPostgres(),
		adapter.WithMessagePublisher(),
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

	a.exportCore = deps.ExportCore
	busConfig := postgresBus.Config{
		PollInterval:      time.Duration(config.Envs.MessageBus.PostgresPollIntervalMS) * time.Millisecond,
		BatchSize:         config.Envs.MessageBus.PostgresBatchSize,
		RetryDelay:        time.Duration(config.Envs.MessageBus.RetryDelaySeconds) * time.Second,
		MaxAttempts:       config.Envs.MessageBus.MaxAttempts,
		ProcessingTimeout: time.Duration(config.Envs.MessageBus.ProcessingTimeoutSecs) * time.Second,
	}

	a.router, a.routerClose, err = postgresBus.NewConsumerRouter(
		adapter.Adapters.Postgres,
		busConfig,
	)
	if err != nil {
		return err
	}

	a.handlers = nil

	userAdapter := userconsumer.NewAdapter()
	if err = a.registerConsumerHandlers(busConfig, userAdapter.HandlerDefinitions()); err != nil {
		return err
	}
	userInvitationAdapter := userinvitationconsumer.NewAdapter()
	if err = a.registerConsumerHandlers(busConfig, userInvitationAdapter.HandlerDefinitions()); err != nil {
		return err
	}
	exportJobAdapter := exportjobconsumer.NewAdapter(exportjobconsumer.Config{
		Core: a.exportCore,
	})
	if err = a.registerConsumerHandlers(busConfig, exportJobAdapter.HandlerDefinitions()); err != nil {
		return err
	}

	if a.shutdown == nil {
		a.shutdown = deps.Shutdown
	}

	logRegisteredHandlers(ctx, a.handlers)

	return nil
}

func logRegisteredHandlers(ctx context.Context, handlers []consumerHandlerInfo) {
	if len(handlers) == 0 {
		log.Ctx(ctx).Warn().Msg("consumer started without registered handlers")
		return
	}

	log.Ctx(ctx).Info().
		Int("handler_count", len(handlers)).
		Msg("consumer handlers registered")

	for _, handler := range handlers {
		log.Ctx(ctx).Info().
			Str("handler_name", handler.Name).
			Str("topic", handler.Topic).
			Str("consumer_group", handler.ConsumerGroup).
			Str("description", handler.Description).
			Msg("consumer handler active")
	}
}
