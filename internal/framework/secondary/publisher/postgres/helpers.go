package postgres

import (
	"time"
)

type Config struct {
	PollInterval      time.Duration
	BatchSize         int
	RetryDelay        time.Duration
	MaxAttempts       int
	ProcessingTimeout time.Duration
}

func (c Config) normalized() Config {
	cfg := c
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = time.Second
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 10
	}
	if cfg.RetryDelay <= 0 {
		cfg.RetryDelay = 30 * time.Second
	}
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 5
	}
	if cfg.ProcessingTimeout <= 0 {
		cfg.ProcessingTimeout = 5 * time.Minute
	}

	return cfg
}
