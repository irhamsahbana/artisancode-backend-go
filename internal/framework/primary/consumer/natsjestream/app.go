package consumer

import (
	"codebase-app/internal/adapter"
	exportJobCore "codebase-app/internal/core/export_job"
	"codebase-app/internal/entity/common"
	attendanceRepo "codebase-app/internal/framework/secondary/db/postgres/attendance"
	exportJobRepo "codebase-app/internal/framework/secondary/db/postgres/export_job"
	storageRepo "codebase-app/internal/framework/secondary/db/postgres/storage"
	natsjetstream "codebase-app/internal/framework/secondary/publisher/natsjetstream"
	"codebase-app/internal/infrastructure/config"
	infraTracing "codebase-app/internal/infrastructure/tracing"
	storage "codebase-app/internal/integration/storage"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {
	ctx, span := infraTracing.StartSpan(ctx, "consumer.Run")
	defer span.End()

	emailConsumer, exportConsumer, exportCore, exportPublisher, exportConsumerManager, err := a.build(ctx)
	if err != nil {
		infraTracing.RecordError(span, err)
		return err
	}
	defer func() {
		if err := exportConsumerManager.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close export consumer manager")
		}
		if err := exportPublisher.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close export message publisher")
		}
	}()

	emailConsumeCtx, err := emailConsumer.Consume(func(msg jetstream.Msg) {
		msgCtx := otel.GetTextMapPropagator().Extract(ctx, messageHeadersCarrier(msg.Headers()))

		switch msg.Subject() {
		case common.MessageSubjectEmailVerification:
			EmailVerificationHandler(msgCtx, msg)
		case common.MessageSubjectEmailForgotPassword:
			ForgotPasswordHandler(msgCtx, msg)
		default:
		}
	})
	if err != nil {
		infraTracing.RecordError(span, err)
		return err
	}
	adapter.Adapters.EmailConsumerCtxNats = emailConsumeCtx

	exportConsumeCtx, err := exportConsumer.Consume(ExportJobRequestedHandler(ctx, exportCore))
	if err != nil {
		emailConsumeCtx.Stop()
		infraTracing.RecordError(span, err)
		return err
	}
	defer exportConsumeCtx.Stop()

	err = a.waitForShutdown(ctx)
	if err != nil {
		infraTracing.RecordError(span, err)
	}

	return err
}

func (a *App) build(ctx context.Context) (jetstream.Consumer, integrationPorts.MessageBusConsumer, integrationExportJobProcessor, integrationPorts.MessagePublisher, integrationPorts.MessageConsumerManager, error) {
	var emailConsumerCtx jetstream.ConsumeContext

	adapter.Adapters.Sync(
		adapter.WithEmailConsumerNats(emailConsumerCtx),
		adapter.WithStorage(),
	)

	exportPublisher, err := natsjetstream.NewPublisher(config.Envs.EmailVerificationQueueNats.NatsURL)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	exportConsumerManager, err := natsjetstream.NewConsumerManager(config.Envs.EmailVerificationQueueNats.NatsURL)
	if err != nil {
		_ = exportPublisher.Close()
		return nil, nil, nil, nil, nil, err
	}

	exportConsumer, err := exportConsumerManager.CreateConsumer(ctx, integrationPorts.MessageBusConsumerConfig{
		StreamName:          common.MessageStreamExportJobService,
		StreamDescription:   "Export job service stream",
		Subjects:            []string{common.MessageSubjectExportJobRequested},
		MaxBytes:            1024 * 1024 * 1024,
		MaxAge:              time.Hour * 24 * 14,
		ConsumerName:        common.MessageConsumerExportJobService,
		Durable:             common.MessageConsumerExportJobService,
		ConsumerDescription: "Export job consumer",
	})
	if err != nil {
		_ = exportConsumerManager.Close()
		_ = exportPublisher.Close()
		return nil, nil, nil, nil, nil, err
	}

	exportJobRepository := exportJobRepo.NewExportJobRepository(exportJobRepo.ExportJobRepositoryConfig{
		DB: adapter.Adapters.Postgres,
	})
	attendanceRepository := attendanceRepo.NewAttendanceRepository(attendanceRepo.AttendanceRepositoryConfig{
		DB: adapter.Adapters.Postgres,
	})
	storageRepository := storageRepo.NewStorageRepository(storageRepo.StorageRepositoryConfig{
		DB: adapter.Adapters.Postgres,
	})
	s3 := storage.NewStorageIntegration(adapter.Adapters.Storage)
	exportCore := exportJobCore.NewExportJobCore(exportJobCore.ExportJobCoreConfig{
		Repo:           exportJobRepository,
		AttendanceRepo: attendanceRepository,
		StorageRepo:    storageRepository,
		S3:             s3,
		Bus:            exportPublisher,
	})

	return adapter.Adapters.EmailConsumerNats, exportConsumer, exportCore, exportPublisher, exportConsumerManager, nil
}

func (a *App) waitForShutdown(ctx context.Context) error {
	ctx, span := infraTracing.StartSpan(ctx, "consumer.WaitForShutdown")
	defer span.End()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Ctx(ctx).Info().Msg("Consumer gracefully stopped")

	err := adapter.Adapters.Unsync()
	if err != nil {
		infraTracing.RecordError(span, err)
	}

	return err
}
