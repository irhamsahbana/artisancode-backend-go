package cmd

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/infrastructure"
	"codebase-app/internal/infrastructure/config"
	infraLogging "codebase-app/internal/infrastructure/logging"
	infraTracing "codebase-app/internal/infrastructure/tracing"
	storage "codebase-app/internal/integration/storage"
	"codebase-app/internal/middleware"
	"codebase-app/internal/setup"
	"codebase-app/pkg/validator"
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/rs/zerolog/log"
)

func RunServer(cmd *flag.FlagSet, args []string) {
	var (
		envs        = config.Envs
		flagAppPort = cmd.String("port", "3000", "Application port")
		SERVER_PORT string
	)

	app := fiber.New()
	adapter.Adapters.Sync(
		adapter.WithRestServer(app),
		adapter.WithOpenAISDK(),
		adapter.WithPostgres(),
		adapter.WithValidator(validator.NewValidator()),
		adapter.WithStorage(),
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
		SERVER_PORT = envs.App.Port
	} else {
		SERVER_PORT = *flagAppPort
	}

	var (
		s3 = storage.NewStorageIntegration(adapter.Adapters.Storage)
	)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS,HEAD",
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
	}))
	app.Use(middleware.RequestID)
	app.Use(middleware.WithAppLogger(log.Logger))
	app.Use(middleware.WithTracing(envs.App.Name))
	app.Use(middleware.Recover())
	app.Use(middleware.WithAccessLog(infrastructure.AccessLogger))
	// End Application Middlewares

	metricTitle := envs.App.Name + " " + envs.App.Environtment + " " + "Metrics"
	app.Get("/metrics", monitor.New(monitor.Config{Title: metricTitle}))

	setup.Dependencies(
		app,
		adapter.Adapters.Postgres,
		s3,
	)

	// Run server in goroutine
	go func() {
		// tell where to connect if in same network
		log.Info().Msgf("Server is running on port %s", SERVER_PORT)
		// Print local addresses for convenience
		log.Info().Msgf("Connect via: http://localhost:%s", SERVER_PORT)
		for _, ip := range getLocalIPv4s() {
			log.Info().Msgf("Connect via: http://%s:%s", ip, SERVER_PORT)
		}
		if err := app.Listen(":" + SERVER_PORT); err != nil {
			log.Fatal().Msgf("Error while starting server: %v", err)
		}
	}()
	// End Run server in goroutine

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)

	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	if runtime.GOOS == "windows" {
		shutdownSignals = []os.Signal{os.Interrupt}
	}

	signal.Notify(quit, shutdownSignals...)
	<-quit
	log.Info().Msg("Server is shutting down ...")

	err = adapter.Adapters.Unsync()
	if err != nil {
		log.Error().Msgf("Error while closing adapters: %v", err)
	}

	log.Info().Msg("Server gracefully stopped")
}

// getLocalIPv4s returns non-loopback IPv4 addresses for network interfaces that are up
func getLocalIPv4s() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	uniq := make(map[string]struct{})
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil { // skip IPv6
				continue
			}
			uniq[ip.String()] = struct{}{}
		}
	}
	out := make([]string, 0, len(uniq))
	for ip := range uniq {
		out = append(out, ip)
	}
	return out
}
