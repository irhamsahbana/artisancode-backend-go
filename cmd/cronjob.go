package cmd

import (
	"codebase-app/internal/adapter"
	exportJobCore "codebase-app/internal/core/export_job"
	messageQueueCore "codebase-app/internal/core/message_queue"
	storageCore "codebase-app/internal/core/storage"
	"codebase-app/internal/entity/coreentity"
	attendanceRepo "codebase-app/internal/framework/secondary/db/postgres/attendance"
	exportJobRepo "codebase-app/internal/framework/secondary/db/postgres/export_job"
	messageQueueRepo "codebase-app/internal/framework/secondary/db/postgres/message_queue"
	storageRepo "codebase-app/internal/framework/secondary/db/postgres/storage"
	storageIntegration "codebase-app/internal/integration/storage"
	"context"
	"flag"
	"os"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
)

// RunCronjob runs the file upload process depending on file size
func RunCronjob(cmd *flag.FlagSet, args []string) {
	var (
		task           = cmd.String("task", "backup", "cron task to run")
		limit          = cmd.Int("limit", 100, "max number of rows/files to clean per run")
		retentionHours = cmd.Int("retention-hours", 168, "retention in hours for processed message queue cleanup")
	)

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing cronjob flags")
	}

	adapter.Adapters.Sync(
		adapter.WithPostgres(),
		adapter.WithStorage(),
	)

	if *task == "cleanup-expired-storage-files" {
		runCleanupExpiredStorageFiles(*limit)
		return
	}
	if *task == "process-export-jobs" {
		runProcessExportJobs(*limit)
		return
	}
	if *task == "cleanup-processed-message-queue" {
		runCleanupProcessedMessageQueue(*limit, *retentionHours)
		return
	}
}

func runCleanupExpiredStorageFiles(limit int) {
	repo := storageRepo.NewStorageRepository(storageRepo.StorageRepositoryConfig{
		DB: adapter.Adapters.Postgres,
	})
	s3 := storageIntegration.NewStorageIntegration(adapter.Adapters.Storage)
	core := storageCore.NewStorageCore(s3, repo)

	resp, err := core.CleanupExpiredFiles(context.Background(), coreentity.CleanupExpiredFilesReq{
		Limit: limit,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to clean up expired storage files")
		return
	}

	log.Info().
		Int("scanned", resp.Scanned).
		Int("deleted", resp.Deleted).
		Msg("Expired storage file cleanup completed")
}

func runProcessExportJobs(limit int) {
	exportJobRepository := exportJobRepo.NewExportJobRepository(exportJobRepo.ExportJobRepositoryConfig{
		DB: adapter.Adapters.Postgres,
	})
	attendanceRepository := attendanceRepo.NewAttendanceRepository(attendanceRepo.AttendanceRepositoryConfig{
		DB: adapter.Adapters.Postgres,
	})
	storageRepository := storageRepo.NewStorageRepository(storageRepo.StorageRepositoryConfig{
		DB: adapter.Adapters.Postgres,
	})
	s3 := storageIntegration.NewStorageIntegration(adapter.Adapters.Storage)
	core := exportJobCore.NewExportJobCore(exportJobCore.ExportJobCoreConfig{
		Repo:           exportJobRepository,
		AttendanceRepo: attendanceRepository,
		StorageRepo:    storageRepository,
		S3:             s3,
	})

	resp, err := core.ProcessPendingExportJobs(context.Background(), limit)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to process export jobs")
		return
	}

	log.Info().
		Int("processed", resp.Processed).
		Msg("Export job processing completed")
}

func runCleanupProcessedMessageQueue(limit, retentionHours int) {
	repo := messageQueueRepo.NewMessageQueueRepository(messageQueueRepo.MessageQueueRepositoryConfig{
		DB: adapter.Adapters.Postgres,
	})
	core := messageQueueCore.NewMessageQueueCore(messageQueueCore.MessageQueueCoreConfig{
		Repo: repo,
	})

	before := time.Now().UTC().Add(-time.Duration(retentionHours) * time.Hour).Format(time.RFC3339)
	resp, err := core.CleanupProcessedMessages(context.Background(), coreentity.CleanupProcessedMessageQueueReq{
		Before: before,
		Limit:  limit,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to clean processed message queue")
		return
	}

	log.Info().
		Int("deleted", resp.Deleted).
		Int("limit", limit).
		Int("retention_hours", retentionHours).
		Msg("Processed message queue cleanup completed")
}

// execCommand runs a command with the given arguments and logs the output
func execCommand(name string, args ...string) error {
	// log command execution
	log.Info().Str("command", name).Strs("args", args).Msg("Executing command")

	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
