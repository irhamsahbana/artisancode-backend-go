package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *meCore) GetMe(ctx context.Context, filter coreentity.SelfFilter) (*coreentity.Me, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetMe")
	defer span.End()

	return &coreentity.Me{
		UserCtx:     filter.UserCtx,
		UserID:      filter.UserCtx.UserID,
		UserName:    filter.UserCtx.UserName,
		TenantID:    filter.UserCtx.TenantID,
		TenantName:  filter.UserCtx.TenantName,
		Roles:       filter.UserCtx.Roles,
		CompanyID:   filter.UserCtx.CompanyID,
		CompanyName: filter.UserCtx.CompanyName,
	}, nil
}
