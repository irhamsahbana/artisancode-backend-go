package cmd

import (
	"codebase-app/internal/adapter"
	httpModule "codebase-app/internal/framework/primary/http"
	"codebase-app/internal/infrastructure/config"
	infraLogging "codebase-app/internal/infrastructure/logging"
	infraTracing "codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/validator"
	"context"
	"flag"

	"github.com/rs/zerolog/log"
)

func RunHttpServer(cmd *flag.FlagSet, args []string) {
	var (
		envs        = config.Envs
		flagAppPort = cmd.String("port", "3000", "Application port")
		serverPort  string
	)

	adapter.Adapters.Sync(
		adapter.WithPostgres(),
		adapter.WithValidator(validator.NewValidator()),
	)

	var otlpEndpoint string
	if envs.Instrumentation.Enabled {
		otlpEndpoint = envs.Instrumentation.OtlpEndpoint
	}

	lp, logWriter, err := infraLogging.InitLogger(&infraLogging.Config{
		Endpoint:      otlpEndpoint,
		AppName:       envs.App.Name,
		AppVersion:    envs.App.Version,
		AppEnv:        envs.App.Environtment,
		LogFile:       envs.App.LogFile,
		AccessLogFile: envs.App.LogFileAccess,
		LogLevel:      envs.App.LogLevel,
		DB:            adapter.Adapters.Postgres,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize logger")
	} else {
		defer func() {
			if err := lp.Shutdown(context.Background()); err != nil {
				log.Error().Err(err).Msg("Error shutting down logger provider")
			}
		}()
	}

	if envs.Instrumentation.Enabled {
		tp, traceErr := infraTracing.InitTracer(&infraTracing.Config{
			Endpoint:   envs.Instrumentation.OtlpEndpoint,
			Headers:    envs.Instrumentation.OtlpHeaders,
			Insecure:   envs.Instrumentation.OtlpInsecure,
			Debug:      envs.Instrumentation.Debug,
			AppName:    envs.App.Name,
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

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing flags")
	}

	if envs.App.Port != "" {
		serverPort = envs.App.Port
	} else {
		serverPort = *flagAppPort
	}

	app := httpModule.NewApp()
	err = app.Run(
		context.Background(),
		envs.App.Name,
		envs.App.Environtment,
		serverPort,
		envs.EmailVerificationQueueNats.NatsURL,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("http::RunHttpServer::Failed to run HTTP app")
	}
}
