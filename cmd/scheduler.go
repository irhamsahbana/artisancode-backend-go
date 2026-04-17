package cmd

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/infrastructure/config"
	infraLogging "codebase-app/internal/infrastructure/logging"
	infraTracing "codebase-app/internal/infrastructure/tracing"
	"context"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
)

func RunScheduler(cmd *flag.FlagSet, args []string) {
	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing scheduler flags")
	}

	envs := config.Envs

	adapter.Adapters.Sync(
		adapter.WithPostgres(),
		adapter.WithStorage(),
	)
	defer func() {
		if err := adapter.Adapters.Unsync(); err != nil {
			log.Error().Err(err).Msg("Failed to close scheduler adapters")
		}
	}()

	var otlpEndpoint string
	if envs.Instrumentation.Enabled {
		otlpEndpoint = envs.Instrumentation.OtlpEndpoint
	}

	lp, logWriter, err := infraLogging.InitLogger(&infraLogging.Config{
		Endpoint:      otlpEndpoint,
		AppName:       envs.App.Name,
		AppVersion:    envs.App.Version,
		AppEnv:        envs.App.Environtment,
		LogFile:       filepath.Join("logs", "scheduler.log"),
		AccessLogFile: envs.App.LogFileAccess,
		LogLevel:      envs.App.LogLevel,
		DB:            adapter.Adapters.Postgres,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize logger")
	}
	defer func() {
		if err := lp.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("Error shutting down logger provider")
		}
	}()

	if envs.Instrumentation.Enabled {
		tp, traceErr := infraTracing.InitTracer(&infraTracing.Config{
			Endpoint:   envs.Instrumentation.OtlpEndpoint,
			Headers:    envs.Instrumentation.OtlpHeaders,
			Insecure:   envs.Instrumentation.OtlpInsecure,
			Debug:      envs.Instrumentation.Debug,
			AppName:    envs.App.Name + "-scheduler",
			AppVersion: envs.App.Version,
			AppEnv:     envs.App.Environtment,
			LogWriter:  logWriter,
		})
		if traceErr != nil {
			log.Error().Err(traceErr).Msg("Failed to initialize tracer")
		} else {
			defer func() {
				if err := tp.Shutdown(context.Background()); err != nil {
					log.Error().Err(err).Msg("Error shutting down tracer provider")
				}
			}()
		}
	}

	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize scheduler")
	}
	defer func() {
		if err := scheduler.Shutdown(); err != nil {
			log.Error().Err(err).Msg("Failed to shutdown scheduler")
		}
	}()

	_, err = scheduler.NewJob(
		gocron.CronJob(envs.Scheduler.StorageCleanupSpec, false),
		gocron.NewTask(runCleanupExpiredStorageFiles, envs.Scheduler.StorageCleanupLimit),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		log.Fatal().Err(err).Str("spec", envs.Scheduler.StorageCleanupSpec).Msg("Failed to register storage cleanup scheduler")
	}

	_, err = scheduler.NewJob(
		gocron.CronJob(envs.Scheduler.MessageQueueCleanupSpec, false),
		gocron.NewTask(
			runCleanupProcessedMessageQueue,
			envs.Scheduler.MessageQueueCleanupLimit,
			envs.Scheduler.MessageQueueRetentionHours,
		),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		log.Fatal().Err(err).Str("spec", envs.Scheduler.MessageQueueCleanupSpec).Msg("Failed to register message queue cleanup scheduler")
	}

	scheduler.Start()

	log.Info().
		Str("storage_cleanup_spec", envs.Scheduler.StorageCleanupSpec).
		Int("storage_cleanup_limit", envs.Scheduler.StorageCleanupLimit).
		Str("message_queue_cleanup_spec", envs.Scheduler.MessageQueueCleanupSpec).
		Int("message_queue_cleanup_limit", envs.Scheduler.MessageQueueCleanupLimit).
		Int("message_queue_retention_hours", envs.Scheduler.MessageQueueRetentionHours).
		Msg("Scheduler is running")

	waitForSchedulerShutdown()
}

func waitForSchedulerShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Scheduler gracefully stopped")
}
