package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r *internalClientRepo) GetInternalClientOwnerPermissions(
	ctx context.Context,
	clientID string,
) (*coreentity.InternalClientOwnerPermissions, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalclient:owner_permissions:GetInternalClientOwnerPermissions",
	)
	defer span.End()

	ownerRoleID, err := r.getOwnerRoleID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	available := make([]coreentity.Permission, 0)
	availableQuery := `
		SELECT id, name, COALESCE(description, '') AS description
		FROM permissions
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY name ASC
	`
	if err := r.db.SelectContext(ctx, &available, r.db.Rebind(availableQuery), clientID); err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"client_id": clientID}).
			Msg("Failed to query internal client permissions")
		return nil, err
	}

	type selectedRow struct {
		ID string `db:"id"`
	}
	selectedRows := make([]selectedRow, 0)
	selectedQuery := `
		SELECT p.id
		FROM permissions p
		INNER JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = ? AND p.tenant_id = ? AND p.deleted_at IS NULL
		ORDER BY p.name ASC
	`
	if err := r.db.SelectContext(ctx, &selectedRows, r.db.Rebind(selectedQuery), ownerRoleID, clientID); err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"client_id": clientID}).
			Msg("Failed to query internal client owner permissions")
		return nil, err
	}

	selectedIDs := make([]string, 0, len(selectedRows))
	for _, row := range selectedRows {
		selectedIDs = append(selectedIDs, row.ID)
	}

	return &coreentity.InternalClientOwnerPermissions{
		ClientID:           clientID,
		Available:          available,
		OwnerPermissionIDs: selectedIDs,
	}, nil
}

func (r *internalClientRepo) UpdateInternalClientOwnerPermissions(
	ctx context.Context,
	data coreentity.InternalClientOwnerPermissionUpdate,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalclient:owner_permissions:UpdateInternalClientOwnerPermissions",
	)
	defer span.End()

	ownerRoleID, err := r.getOwnerRoleID(ctx, data.ClientID)
	if err != nil {
		return err
	}

	if len(data.PermissionIDs) > 0 {
		var validCount int
		countQuery := `
			SELECT COUNT(DISTINCT id)
			FROM permissions
			WHERE tenant_id = ? AND deleted_at IS NULL AND id IN (?)
		`
		query, args, err := sqlx.In(countQuery, data.ClientID, data.PermissionIDs)
		if err != nil {
			return err
		}
		query = r.db.Rebind(query)
		if err := r.db.GetContext(ctx, &validCount, query, args...); err != nil {
			log.Ctx(ctx).
				Error().
				Err(err).
				Any(common.LogKeyPayload, data).
				Msg("Failed to validate internal client owner permissions")
			return err
		}
		if validCount != len(data.PermissionIDs) {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Invalid internal client owner permissions")
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageOneOrMorePermissionsAreInvalidForThisClient)
		}
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, data).
			Msg("Failed to begin internal client owner permission transaction")
		return err
	}
	defer tx.Rollback()

	deleteQuery := `DELETE FROM role_permissions WHERE role_id = ? AND tenant_id = ?`
	if _, err := tx.ExecContext(ctx, r.db.Rebind(deleteQuery), ownerRoleID, data.ClientID); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to delete existing owner permissions")
		return err
	}

	if len(data.PermissionIDs) > 0 {
		insertQuery := `INSERT INTO role_permissions (role_id, permission_id, tenant_id) VALUES (?, ?, ?)`
		for _, permissionID := range data.PermissionIDs {
			if _, err := tx.ExecContext(ctx, r.db.Rebind(insertQuery), ownerRoleID, permissionID, data.ClientID); err != nil {
				log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to insert owner permission")
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, data).
			Msg("Failed to commit internal client owner permission transaction")
		return err
	}

	return nil
}

func (r *internalClientRepo) getOwnerRoleID(ctx context.Context, clientID string) (string, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalclient:owner_permissions:getOwnerRoleID",
	)
	defer span.End()

	var ownerRoleID string
	query := `
		SELECT r.id
		FROM roles r
		INNER JOIN tenants t ON t.id = r.tenant_id AND t.deleted_at IS NULL
		WHERE r.tenant_id = ? AND r.name = 'owner' AND r.deleted_at IS NULL
		LIMIT 1
	`
	if err := r.db.GetContext(ctx, &ownerRoleID, r.db.Rebind(query), clientID); err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).
				Warn().
				Any(common.LogKeyPayload, map[string]string{"client_id": clientID}).
				Msg(errmsg.MessageClientOwnerRoleNotFound)
			return "", errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageClientOwnerRoleNotFound)
		}
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"client_id": clientID}).
			Msg("Failed to get internal client owner role")
		return "", err
	}
	return ownerRoleID, nil
}
