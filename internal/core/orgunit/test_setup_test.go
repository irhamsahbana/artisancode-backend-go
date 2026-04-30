package core

import (
	"os"
	"testing"

	"codebase-app/internal/infrastructure/config"
)

func TestMain(m *testing.M) {
	config.Envs = &config.Config{}
	config.Envs.App.Name = "test"
	os.Exit(m.Run())
}
