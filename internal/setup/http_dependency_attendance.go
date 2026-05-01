package setup

import (
	attendanceCore "codebase-app/internal/core/attendance"
	exportJobCore "codebase-app/internal/core/export_job"
	storageCore "codebase-app/internal/core/storage"
	attendanceHandler "codebase-app/internal/framework/primary/http/attendance"
	exportJobHandler "codebase-app/internal/framework/primary/http/export_job"
	storageHandler "codebase-app/internal/framework/primary/http/storage"
	attendanceRepo "codebase-app/internal/framework/secondary/db/postgres/attendance"
	exportJobRepo "codebase-app/internal/framework/secondary/db/postgres/export_job"
	storageRepo "codebase-app/internal/framework/secondary/db/postgres/storage"
	"codebase-app/internal/middleware"
)

func buildAttendanceDependencies(deps *httpDependencies, ctx httpBootstrapContext) {
	attendanceRepository := attendanceRepo.NewAttendanceRepository(attendanceRepo.Config{
		DB: ctx.db,
	})
	exportJobRepository := exportJobRepo.NewExportJobRepository(exportJobRepo.Config{
		DB: ctx.db,
	})
	storageRepository := storageRepo.NewStorageRepository(storageRepo.Config{
		DB: ctx.db,
	})

	deps.attendanceCore = attendanceCore.NewAttendanceCore(attendanceCore.Config{
		Repo:        attendanceRepository,
		StorageRepo: storageRepository,
		Tx:          ctx.tx,
		S3:          ctx.s3,
	})
	deps.exportJobCore = exportJobCore.NewExportJobCore(exportJobCore.Config{
		Repo:           exportJobRepository,
		AttendanceRepo: attendanceRepository,
		StorageRepo:    storageRepository,
		Tx:             ctx.tx,
		S3:             ctx.s3,
		Bus:            ctx.bus,
	})
	deps.storageCore = storageCore.NewStorageCore(ctx.s3, storageRepository)
}

func (deps httpDependencies) registerAttendanceRoutes() {
	attendanceHandlerInst := attendanceHandler.NewAttendanceHandler(
		attendanceHandler.Config{
			Core: deps.attendanceCore,
		},
	)
	attendanceHandlerInst.Register(
		deps.app.Group("/attendance-logs", middleware.Auth),
	)
	attendanceHandlerInst.RegisterSummary(
		deps.app.Group("/attendance-summary", middleware.Auth),
	)
	attendanceHandlerInst.RegisterPolicy(
		deps.app.Group("/attendance-policy", middleware.Auth),
	)

	exportJobHandler.NewExportJobHandler(exportJobHandler.Config{
		Core: deps.exportJobCore,
	}).Register(deps.app.Group("/export-jobs", middleware.Auth))
}

func (deps httpDependencies) registerStorageRoutes() {
	storageHandler.NewStorageHandler(deps.storageCore).Register(
		deps.app.Group("/storage"),
	)
}
