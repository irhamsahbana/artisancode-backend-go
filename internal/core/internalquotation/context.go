package core

import (
	"context"

	"codebase-app/internal/entity/common"
)

func coreentityTenant(ctx context.Context) string {
	return common.GetUserContext(ctx).TenantID
}

func userTenant(ctx context.Context) string {
	return coreentityTenant(ctx)
}
