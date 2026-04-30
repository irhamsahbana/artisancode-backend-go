package core

import (
	"context"
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

func internalInvoiceContext() context.Context {
	return context.WithValue(
		context.Background(),
		common.UserContextKeyClaims,
		common.UserContext{TenantID: "tenant-1", UserID: "user-1"},
	)
}
