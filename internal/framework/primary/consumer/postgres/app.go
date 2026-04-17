package consumer

import (
	"codebase-app/internal/adapter"
	exportJobCore "codebase-app/internal/core/export_job"
	"codebase-app/internal/entity/common"
	attendanceRepo "codebase-app/internal/framework/secondary/db/postgres/attendance"
	exportJobRepo "codebase-app/internal/framework/secondary/db/postgres/export_job"
	storageRepo "codebase-app/internal/framework/secondary/db/postgres/storage"
	messagebus "codebase-app/internal/framework/secondary/publisher/messagebus"
	infraTracing "codebase-app/internal/infrastructure/tracing"
	storage "codebase-app/internal/integration/storage"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
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

	emailConsumer, exportConsumer, exportCore, exportPublisher, consumerManager, err := a.build(ctx)
	if err != nil {
		infraTracing.RecordError(span, err)
		return err
	}
	defer func() {
		if err := exportPublisher.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close export message publisher")
		}
		if err := consumerManager.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close consumer manager")
		}
	}()

	emailConsumeCtx, err := emailConsumer.Consume(func(msg integrationPorts.MessageBusMessage) {
		msgCtx := otel.GetTextMapPropagator().Extract(ctx, MessageHeadersCarrier(msg.Headers()))

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
	defer emailConsumeCtx.Stop()

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

func (a *App) build(ctx context.Context) (integrationPorts.MessageBusConsumer, integrationPorts.MessageBusConsumer, IntegrationExportJobProcessor, integrationPorts.MessagePublisher, integrationPorts.MessageConsumerManager, error) {
	adapter.Adapters.Sync(
		adapter.WithStorage(),
	)

	exportPublisher, err := messagebus.NewPublisher(adapterPostgres())
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	consumerManager, err := messagebus.NewConsumerManager(adapterPostgres())
	if err != nil {
		_ = exportPublisher.Close()
		return nil, nil, nil, nil, nil, err
	}

	emailConsumer, err := consumerManager.CreateConsumer(ctx, integrationPorts.MessageBusConsumerConfig{
		StreamName:          common.MessageStreamEmailService,
		StreamDescription:   "Email service stream",
		Subjects:            []string{common.MessageSubjectEmailAll},
		MaxBytes:            1024 * 1024 * 1024,
		MaxAge:              time.Hour * 24 * 14,
		ConsumerName:        common.MessageConsumerEmailService,
		Durable:             common.MessageConsumerEmailService,
		ConsumerDescription: "Email service consumer",
	})
	if err != nil {
		_ = consumerManager.Close()
		_ = exportPublisher.Close()
		return nil, nil, nil, nil, nil, err
	}

	exportConsumer, err := consumerManager.CreateConsumer(ctx, integrationPorts.MessageBusConsumerConfig{
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
		_ = consumerManager.Close()
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
		DB: adapterPostgres(),
	})
	s3 := storage.NewStorageIntegration(adapter.Adapters.Storage)
	exportCore := exportJobCore.NewExportJobCore(exportJobCore.ExportJobCoreConfig{
		Repo:           exportJobRepository,
		AttendanceRepo: attendanceRepository,
		StorageRepo:    storageRepository,
		S3:             s3,
		Bus:            exportPublisher,
	})

	return emailConsumer, exportConsumer, exportCore, exportPublisher, consumerManager, nil
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

func adapterPostgres() *sqlx.DB {
	return adapter.Adapters.Postgres
}
