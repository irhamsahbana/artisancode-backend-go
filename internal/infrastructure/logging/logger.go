package logging

import (
	"io"

	infra "codebase-app/internal/infrastructure"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Endpoint      string
	AppName       string
	AppVersion    string
	AppEnv        string
	LogFile       string
	AccessLogFile string
	LogLevel      string
	DB            *sqlx.DB
}

type Provider = infra.LoggerProvider

func InitLogger(cfg *Config) (Provider, io.Writer, error) {
	return infra.InitLogger(&infra.LoggerConfig{
		Endpoint:      cfg.Endpoint,
		AppName:       cfg.AppName,
		AppVersion:    cfg.AppVersion,
		AppEnv:        cfg.AppEnv,
		LogFile:       cfg.LogFile,
		AccessLogFile: cfg.AccessLogFile,
		LogLevel:      cfg.LogLevel,
		DB:            cfg.DB,
	})
}
