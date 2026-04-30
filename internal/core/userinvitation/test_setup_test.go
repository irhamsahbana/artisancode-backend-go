package core

import (
	"os"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/config"
)

func TestMain(m *testing.M) {
	config.Envs = &config.Config{}
	config.Envs.App.Name = "test"
	os.Exit(m.Run())
}

func invitationUserContext(roles ...string) common.UserContext {
	return common.UserContext{
		UserID:   "user-1",
		TenantID: "tenant-1",
		Roles:    roles,
	}
}
