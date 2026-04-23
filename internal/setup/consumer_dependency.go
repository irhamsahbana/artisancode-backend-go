package setup

import (
	"context"
	"time"

	"codebase-app/internal/adapter"
	exportJobCore "codebase-app/internal/core/export_job"
	"codebase-app/internal/entity/common"
	attendanceRepo "codebase-app/internal/framework/secondary/db/postgres/attendance"
	exportJobRepo "codebase-app/internal/framework/secondary/db/postgres/export_job"
	storageRepo "codebase-app/internal/framework/secondary/db/postgres/storage"
	storage "codebase-app/internal/integration/storage"
	corePorts "codebase-app/internal/ports/core"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
)

type ConsumerDependencies struct {
	EmailSubscription  integrationPorts.MessageBusSubscription
	ExportSubscription integrationPorts.MessageBusSubscription
	ExportCore         corePorts.ExportJobCore
	Shutdown           func() error
}

func NewConsumerDependencies(
	ctx context.Context,
	shutdown func() error,
) (ConsumerDependencies, error) {
	var (
		db                  = adapter.Adapters.Postgres
		bus                 = adapter.Adapters.MessagePublisher
		s3                  = storage.NewStorageIntegration(adapter.Adapters.Storage)
		subscriptionManager = adapter.Adapters.MessageSubscriptionManager
	)

	emailSubscription, err := subscriptionManager.CreateSubscription(ctx, common.MessageBusSubscriptionConfig{
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
		return ConsumerDependencies{}, err
	}

	exportSubscription, err := subscriptionManager.CreateSubscription(ctx, common.MessageBusSubscriptionConfig{
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
		return ConsumerDependencies{}, err
	}

	exportJobRepository := exportJobRepo.NewExportJobRepository(exportJobRepo.ExportJobRepositoryConfig{
		DB: db,
	})
	attendanceRepository := attendanceRepo.NewAttendanceRepository(attendanceRepo.AttendanceRepositoryConfig{
		DB: db,
	})
	storageRepository := storageRepo.NewStorageRepository(storageRepo.StorageRepositoryConfig{
		DB: db,
	})
	exportCore := exportJobCore.NewExportJobCore(exportJobCore.ExportJobCoreConfig{
		Repo:           exportJobRepository,
		AttendanceRepo: attendanceRepository,
		StorageRepo:    storageRepository,
		S3:             s3,
		Bus:            bus,
	})

	return ConsumerDependencies{
		EmailSubscription:  emailSubscription,
		ExportSubscription: exportSubscription,
		ExportCore:         exportCore,
		Shutdown:           shutdown,
	}, nil
}
