package infrastructure

import (
	"context"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/jmoiron/sqlx"
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
	DB            *sqlx.DB
}

type dbLogWriter struct {
	db      *sqlx.DB
	logChan chan string
}

func newDBLogWriter(db *sqlx.DB) *dbLogWriter {
	w := &dbLogWriter{
		db:      db,
		logChan: make(chan string, 1000), // buffered channel to prevent blocking
	}
	// Start background goroutine to process logs
	go w.processLogs()
	return w
}

func (w *dbLogWriter) processLogs() {
	for log := range w.logChan {
		if w.db != nil {
			_, _ = w.db.ExecContext(context.Background(), "INSERT INTO logs (log) VALUES ($1::jsonb)", log)
		}
	}
}

func (w *dbLogWriter) Write(p []byte) (n int, err error) {
	if w.db == nil {
		return len(p), nil
	}
	raw := strings.TrimSpace(string(p))
	if raw == "" {
		return len(p), nil
	}
	// Non-blocking send to channel
	select {
	case w.logChan <- raw:
		// Successfully sent to channel
	default:
		// Channel full, skip this log to prevent blocking
		// Could add a counter here to track dropped logs
	}
	return len(p), nil
}

func (w *dbLogWriter) Close() error {
	if w.logChan != nil {
		close(w.logChan)
	}
	return nil
}

func InitLogger(cfg *LoggerConfig) (LoggerProvider, io.Writer, error) {
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	InitializeLogger(cfg.AppEnv, cfg.LogFile, logLevel, cfg.DB)
	InitializeAccessLogger(cfg.AppEnv, cfg.AccessLogFile, logLevel, cfg.DB)

	return &loggerProviderNoop{}, os.Stderr, nil
}

// InitializeLogger will set logging format.
func InitializeLogger(stage string, filename string, logLevel zerolog.Level, db *sqlx.DB) {
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
	if db != nil {
		writers = append(writers, newDBLogWriter(db))
	}
	mw := io.MultiWriter(writers...)

	// using json format for production
	var logger zerolog.Logger
	if stage == "production" {
		if db != nil {
			logger = zerolog.New(io.MultiWriter(lumberjackLogger, newDBLogWriter(db))).With().Timestamp().Caller().Logger().Level(zerolog.InfoLevel)
		} else {
			logger = zerolog.New(lumberjackLogger).With().Timestamp().Caller().Logger().Level(zerolog.InfoLevel)
		}
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
			<-q
			lumberjackLogger.Close()
			log.Info().Msg("Closing logs ...")
			// telegramHook.Stop()
			log.Info().Msg("Closing telegram hook ...")
		}
	}()
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
func InitializeAccessLogger(stage string, filename string, logLevel zerolog.Level, db *sqlx.DB) {
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
	if db != nil {
		writers = append(writers, newDBLogWriter(db))
	}
	mw := io.MultiWriter(writers...)

	var logger zerolog.Logger
	if stage == "production" {
		if db != nil {
			logger = zerolog.New(io.MultiWriter(lumberjackLogger, newDBLogWriter(db))).With().Timestamp().Caller().Logger().Level(zerolog.InfoLevel)
		} else {
			logger = zerolog.New(lumberjackLogger).With().Timestamp().Caller().Logger().Level(zerolog.InfoLevel)
		}
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
