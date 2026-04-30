package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *userCore) GetTenantProfile(ctx context.Context) (*coreentity.TenantProfile, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:get_tenant_profile:GetTenantProfile")
	defer span.End()

	userCtx := common.GetUserContext(ctx)
	if userCtx.TenantID == "" {
		log.Ctx(ctx).Warn().Msg("Tenant profile requested without tenant context")
		return nil, errmsg.NewCustomErrors(401).SetMessage("Invalid credentials")
	}

	return c.repo.GetTenantProfile(ctx, userCtx.TenantID)
}
