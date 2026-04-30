package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"

	"github.com/stretchr/testify/require"
)

func TestGetMe_ReturnsUserContextProfile(t *testing.T) {
	companyID := "company-1"
	companyName := "Artisan"
	userCtx := common.UserContext{
		UserID:      "user-1",
		UserName:    "Budi",
		TenantID:    "tenant-1",
		TenantName:  "Tenant",
		Roles:       []string{"employee"},
		CompanyID:   &companyID,
		CompanyName: &companyName,
	}
	core := NewMeCore(Config{})

	got, err := core.GetMe(context.Background(), coreentity.SelfFilter{UserCtx: userCtx})

	require.NoError(t, err)
	require.Equal(t, &coreentity.Me{
		UserCtx:     userCtx,
		UserID:      "user-1",
		UserName:    "Budi",
		TenantID:    "tenant-1",
		TenantName:  "Tenant",
		Roles:       []string{"employee"},
		CompanyID:   &companyID,
		CompanyName: &companyName,
	}, got)
}
