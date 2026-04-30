package core

import (
	"codebase-app/internal/infrastructure/config"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config.Envs = &config.Config{}
	config.Envs.App.Name = "test"
	os.Exit(m.Run())
}
