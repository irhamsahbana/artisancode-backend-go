package setup

import (
	"context"

	"codebase-app/internal/adapter"
	exportJobCore "codebase-app/internal/core/export_job"
	attendanceRepo "codebase-app/internal/framework/secondary/db/postgres/attendance"
	exportJobRepo "codebase-app/internal/framework/secondary/db/postgres/export_job"
	storageRepo "codebase-app/internal/framework/secondary/db/postgres/storage"
	postgresTx "codebase-app/internal/framework/secondary/db/postgres/transaction"
	storage "codebase-app/internal/integration/storage"
	corePorts "codebase-app/internal/ports/core"
)

type ConsumerDependencies struct {
	ExportCore corePorts.ExportJobCore
	Shutdown   func() error
}

func NewConsumerDependencies(
	ctx context.Context,
	shutdown func() error,
) (ConsumerDependencies, error) {
	var (
		db  = adapter.Adapters.Postgres
		bus = adapter.Adapters.MessagePublisher
		s3  = storage.NewStorageIntegration(adapter.Adapters.Storage)
		tx  = postgresTx.NewTransactor(db)
	)

	exportJobRepository := exportJobRepo.NewExportJobRepository(exportJobRepo.Config{
		DB: db,
	})
	attendanceRepository := attendanceRepo.NewAttendanceRepository(attendanceRepo.Config{
		DB: db,
	})
	storageRepository := storageRepo.NewStorageRepository(storageRepo.Config{
		DB: db,
	})
	exportCore := exportJobCore.NewExportJobCore(exportJobCore.Config{
		Repo:           exportJobRepository,
		AttendanceRepo: attendanceRepository,
		StorageRepo:    storageRepository,
		Tx:             tx,
		S3:             s3,
		Bus:            bus,
	})

	return ConsumerDependencies{
		ExportCore: exportCore,
		Shutdown:   shutdown,
	}, nil
}
