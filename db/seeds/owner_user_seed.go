package seeds

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultOwnerName     = "System Owner"
	defaultOwnerUserName = "owner"
	defaultOwnerEmail    = "owner@default.tenant.com"
	defaultOwnerPassword = "owner12345"
)

func (s *Seed) ownerUserSeed() {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	tenantID, err := upsertTenant(tx, defaultTenantName)
	if err != nil {
		return
	}

	var globalTenantID *string
	ownerRoleID, err := getRoleID(tx, globalTenantID, "owner")
	if err != nil {
		return
	}

	if err := upsertOwnerUser(tx, tenantID, ownerRoleID); err != nil {
		return
	}

	log.Info().Msg("owner user seeded successfully!")
}

func getRoleID(tx *sqlx.Tx, tenantID *string, roleName string) (string, error) {
	var roleID string
	query := `
		SELECT id
		FROM roles
		WHERE name = ? AND tenant_id IS NOT DISTINCT FROM ? AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT 1
	`
	if err := tx.Get(&roleID, tx.Rebind(query), roleName, tenantID); err != nil {
		log.Error().Err(err).Any("role_name", roleName).Msg("failed to get role id")
		return "", err
	}
	return roleID, nil
}

func upsertOwnerUser(tx *sqlx.Tx, tenantID, roleID string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultOwnerPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash owner password")
		return err
	}

	var userID string
	query := `
		INSERT INTO users (tenant_id, company_id, name, username, email, password)
		VALUES (?, NULL, ?, ?, ?, ?)
		ON CONFLICT (email) DO UPDATE SET
			tenant_id = EXCLUDED.tenant_id,
			company_id = EXCLUDED.company_id,
			name = EXCLUDED.name,
			username = EXCLUDED.username,
			password = EXCLUDED.password,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`
	if err := tx.Get(&userID, tx.Rebind(query), tenantID, defaultOwnerName, defaultOwnerUserName, defaultOwnerEmail, string(hashedPassword)); err != nil {
		log.Error().Err(err).Msg("failed to upsert owner user")
		return err
	}

	roleQuery := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES (?, ?)
		ON CONFLICT DO NOTHING
	`
	if _, err := tx.Exec(tx.Rebind(roleQuery), userID, roleID); err != nil {
		log.Error().Err(err).Msg("failed to insert user_role")
		return err
	}

	return nil
}
