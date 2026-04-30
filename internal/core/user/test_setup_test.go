package core

import (
	"os"
	"testing"
	"time"

	"codebase-app/internal/infrastructure/config"
)

func TestMain(m *testing.M) {
	config.Envs = &config.Config{}
	config.Envs.App.Name = "test"
	os.Exit(m.Run())
}

func strPtr(s string) *string {
	return &s
}

func timePtr(s string) *time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return &t
}
