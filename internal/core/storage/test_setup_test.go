package core

import infraConfig "codebase-app/internal/infrastructure/config"

func init() {
	infraConfig.Envs = &infraConfig.Config{}
	infraConfig.Envs.App.Name = "storage-core-test"
}

func strPtr(s string) *string { return &s }
