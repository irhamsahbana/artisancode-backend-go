package cmd

import (
	"codebase-app/internal/adapter"
	postgresConsumerModule "codebase-app/internal/framework/primary/consumer/postgres"
	"codebase-app/internal/infrastructure/config"
	infraLogging "codebase-app/internal/infrastructure/logging"
	infraTracing "codebase-app/internal/infrastructure/tracing"
	"context"
	"flag"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func RunConsumer(cmd *flag.FlagSet, args []string) {
	envs := config.Envs

	var otlpEndpoint string
	if envs.Instrumentation.Enabled {
		otlpEndpoint = envs.Instrumentation.OtlpEndpoint
	}

	lp, logWriter, err := infraLogging.InitLogger(&infraLogging.Config{
		Endpoint:      otlpEndpoint,
		AppName:       envs.App.Name,
		AppVersion:    envs.App.Version,
		AppEnv:        envs.App.Environtment,
		LogFile:       filepath.Join("logs", "consumer.log"),
		AccessLogFile: envs.App.LogFileAccess,
		LogLevel:      envs.App.LogLevel,
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
			AppName:    envs.App.Name + "-consumer",
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

	log.Info().Msg("Running consumer")

	app := postgresConsumerModule.NewApp(postgresConsumerModule.AppConfig{
		Shutdown: adapter.Adapters.Unsync,
	})

	err = app.Run(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run consumer app")
	}
}
