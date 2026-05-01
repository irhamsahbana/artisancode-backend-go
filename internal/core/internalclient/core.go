package core

import (
	"context"
	"strings"

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

type Config struct {
	Repo portsRepo.InternalClientRepository
}

var _ corePorts.InternalClientCore = &internalClientCore{}

func NewInternalClientCore(cfg Config) *internalClientCore {
	return &internalClientCore{repo: cfg.Repo}
}

func (c *internalClientCore) GetInternalClients(
	ctx context.Context,
	filter coreentity.InternalClientListFilter,
) ([]coreentity.InternalClient, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalclient:core:GetInternalClients")
	defer span.End()

	if err := c.authorizeRead(common.GetUserContext(ctx)); err != nil {
		return nil, 0, err
	}

	return c.repo.GetInternalClients(ctx, filter)
}

func (c *internalClientCore) GetInternalClientOwnerPermissions(
	ctx context.Context,
	clientID string,
) (*coreentity.InternalClientOwnerPermissions, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalclient:core:GetInternalClientOwnerPermissions")
	defer span.End()

	if err := c.authorizeRead(common.GetUserContext(ctx)); err != nil {
		return nil, err
	}

	return c.repo.GetInternalClientOwnerPermissions(ctx, clientID)
}

func (c *internalClientCore) UpdateInternalClientOwnerPermissions(
	ctx context.Context,
	data coreentity.InternalClientOwnerPermissionUpdate,
) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalclient:core:UpdateInternalClientOwnerPermissions")
	defer span.End()

	if err := c.authorizeManage(common.GetUserContext(ctx)); err != nil {
		return err
	}

	data.PermissionIDs = uniqueNonEmptyStrings(data.PermissionIDs)
	return c.repo.UpdateInternalClientOwnerPermissions(ctx, data)
}

func (c *internalClientCore) authorizeRead(userCtx common.UserContext) error {
	if !userCtx.HasRole(coreentity.InternalUserRoleSuperAdmin) &&
		!userCtx.HasRole(coreentity.InternalUserRoleOperator) {
		return errmsg.NewCustomErrors(403).SetMessage(errmsg.MessageYouAreNotAuthorizedToViewInternalClients)
	}

	return nil
}

func (c *internalClientCore) authorizeManage(userCtx common.UserContext) error {
	if !userCtx.HasRole(coreentity.InternalUserRoleSuperAdmin) &&
		!userCtx.HasRole(coreentity.InternalUserRoleOperator) {
		return errmsg.NewCustomErrors(403).SetMessage(errmsg.MessageYouAreNotAuthorizedToManageInternalClientPermissions)
	}

	return nil
}

func uniqueNonEmptyStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}
