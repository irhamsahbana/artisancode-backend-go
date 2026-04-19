package setup

import (
	"context"
	"time"

	exportJobCore "codebase-app/internal/core/export_job"
	"codebase-app/internal/entity/common"
	consumerApp "codebase-app/internal/framework/primary/consumer/postgres"
	attendanceRepo "codebase-app/internal/framework/secondary/db/postgres/attendance"
	exportJobRepo "codebase-app/internal/framework/secondary/db/postgres/export_job"
	storageRepo "codebase-app/internal/framework/secondary/db/postgres/storage"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/jmoiron/sqlx"
)

func NewConsumerAppConfig(
	ctx context.Context,
	db *sqlx.DB,
	s3 integrationPorts.StorageContract,
	bus integrationPorts.MessagePublisher,
	subscriptionManager integrationPorts.MessageSubscriptionManager,
	shutdown func() error,
) (consumerApp.AppConfig, error) {
	emailSubscription, err := subscriptionManager.CreateSubscription(ctx, integrationPorts.MessageBusSubscriptionConfig{
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
		return consumerApp.AppConfig{}, err
	}

	exportSubscription, err := subscriptionManager.CreateSubscription(ctx, integrationPorts.MessageBusSubscriptionConfig{
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
		return consumerApp.AppConfig{}, err
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

	return consumerApp.AppConfig{
		EmailSubscription:   emailSubscription,
		ExportSubscription:  exportSubscription,
		ExportCore:          exportCore,
		ExportPublisher:     bus,
		SubscriptionManager: subscriptionManager,
		Shutdown:            shutdown,
	}, nil
}
