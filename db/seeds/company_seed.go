package seeds

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const (
	defaultTenantName = "Default Tenant"
	defaultTenantCode = "DEFAULT"
)

func (s *Seed) tenantSeed() {
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

	_, err = upsertTenant(tx, defaultTenantName)
	if err != nil {
		return
	}

	log.Info().Msg("tenant seeded successfully!")
}

func upsertTenant(tx *sqlx.Tx, name string) (string, error) {
	return upsertTenantWithCode(tx, name, defaultTenantCode)
}

func upsertTenantWithCode(tx *sqlx.Tx, name string, code string) (string, error) {
	var id string
	selectQuery := `
		SELECT id
		FROM tenants
		WHERE code = ? AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT 1
	`
	err := tx.Get(&id, tx.Rebind(selectQuery), code)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		log.Error().Err(err).Any("code", code).Msg("failed to query tenant")
		return "", err
	}

	insertQuery := `
		INSERT INTO tenants (name, code)
		VALUES (?, ?)
		RETURNING id
	`
	if err := tx.Get(&id, tx.Rebind(insertQuery), name, code); err != nil {
		log.Error().Err(err).Any("name", name).Msg("failed to insert tenant")
		return "", err
	}
	return id, nil
}
