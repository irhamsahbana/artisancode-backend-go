package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/pkg/errmsg"
)

type internalClientCore struct {
	repo portsRepo.InternalClientRepository
}

type InternalClientCoreConfig struct {
	Repo portsRepo.InternalClientRepository
}

var _ corePorts.InternalClientCore = &internalClientCore{}

func NewInternalClientCore(cfg InternalClientCoreConfig) *internalClientCore {
	return &internalClientCore{repo: cfg.Repo}
}

func (c *internalClientCore) GetInternalClients(ctx context.Context, filter coreentity.InternalClientListFilter) ([]coreentity.InternalClient, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalclient:core:GetInternalClients")
	defer span.End()

	if err := c.authorizeRead(common.GetUserContext(ctx)); err != nil {
		return nil, 0, err
	}

	return c.repo.GetInternalClients(ctx, filter)
}

func (c *internalClientCore) authorizeRead(userCtx common.UserContext) error {
	if !userCtx.HasRole(coreentity.InternalUserRoleSuperAdmin) && !userCtx.HasRole(coreentity.InternalUserRoleOperator) {
		return errmsg.NewCustomErrors(403).SetMessage("You are not authorized to view internal clients")
	}

	return nil
}
