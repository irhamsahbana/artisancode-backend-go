package logging

import (
	"io"

	infra "codebase-app/internal/infrastructure"
)

type Config struct {
	Endpoint      string
	AppName       string
	AppVersion    string
	AppEnv        string
	LogFile       string
	AccessLogFile string
	LogLevel      string
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
	})
}
