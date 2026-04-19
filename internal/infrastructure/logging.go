package infrastructure

import (
	"context"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// var errSkipEvent = errors.New("skip")

// AccessLogger is a dedicated logger for access logs.
var AccessLogger zerolog.Logger

type LoggerProvider interface {
	Shutdown(ctx context.Context) error
}

type loggerProviderNoop struct{}

func (l *loggerProviderNoop) Shutdown(_ context.Context) error {
	return nil
}

type LoggerConfig struct {
	Endpoint      string
	AppName       string
	AppVersion    string
	AppEnv        string
	LogFile       string
	AccessLogFile string
	LogLevel      string
}

func InitLogger(cfg *LoggerConfig) (LoggerProvider, io.Writer, error) {
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	InitializeLogger(cfg.AppEnv, cfg.LogFile, logLevel)
	InitializeAccessLogger(cfg.AppEnv, cfg.AccessLogFile, logLevel)

	return &loggerProviderNoop{}, os.Stderr, nil
}

// InitializeLogger will set logging format.
func InitializeLogger(stage string, filename string, logLevel zerolog.Level) {
	// pr, pw := io.Pipe()
	// ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		log.Fatal().Err(err).Msg("failed to create log directory")
	}

	var (
		lumberjackLogger = &lumberjack.Logger{
			MaxSize:  100,
			MaxAge:   1,
			Filename: filename,
		}
		writers = []io.Writer{zerolog.ConsoleWriter{Out: os.Stderr}, lumberjackLogger}
	)
	mw := io.MultiWriter(writers...)

	// using json format for production
	var logger zerolog.Logger
	if stage == "production" {
		logger = zerolog.New(lumberjackLogger).With().Timestamp().Caller().Logger().Level(zerolog.InfoLevel)
	} else {
		logger = zerolog.New(mw).With().Timestamp().Caller().Logger().Level(logLevel)
	}
	log.Logger = logger

	q := make(chan os.Signal, 1)
	c := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for {
			<-c
			if err := lumberjackLogger.Rotate(); err != nil {
				log.Error().Err(err).Msg("Error while rotating logs")
			}
			log.Info().Msg("Rotating logs ...")
		}
	}()
}

// InitializeAccessLogger configures a separate logger for access logs.
func InitializeAccessLogger(stage string, filename string, logLevel zerolog.Level) {
	// ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		log.Fatal().Err(err).Msg("failed to create access log directory")
	}

	var (
		lumberjackLogger = &lumberjack.Logger{
			MaxSize:  100,
			MaxAge:   1,
			Filename: filename,
		}
		writers = []io.Writer{zerolog.ConsoleWriter{Out: os.Stderr}, lumberjackLogger}
	)
	mw := io.MultiWriter(writers...)

	var logger zerolog.Logger
	if stage == "production" {
		logger = zerolog.New(lumberjackLogger).With().Timestamp().Caller().Logger().Level(zerolog.InfoLevel)
	} else {
		logger = zerolog.New(mw).With().Timestamp().Caller().Logger().Level(logLevel)
	}
	AccessLogger = logger

	q := make(chan os.Signal, 1)
	c := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for {
			<-q
			lumberjackLogger.Close()
			log.Info().Msg("Closing access logs ...")
		}
	}()
	go func() {
		for {
			<-c
			if err := lumberjackLogger.Rotate(); err != nil {
				log.Error().Err(err).Msg("Error while rotating access logs")
			}
			log.Info().Msg("Rotating access logs ...")
		}
	}()
}
