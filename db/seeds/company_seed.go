package seeds

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const defaultTenantName = "Default Tenant"

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
	var id string
	selectQuery := `
		SELECT id
		FROM tenants
		WHERE name = ? AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT 1
	`
	err := tx.Get(&id, tx.Rebind(selectQuery), name)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		log.Error().Err(err).Any("name", name).Msg("failed to query tenant")
		return "", err
	}

	insertQuery := `
		INSERT INTO tenants (name)
		VALUES (?)
		RETURNING id
	`
	if err := tx.Get(&id, tx.Rebind(insertQuery), name); err != nil {
		log.Error().Err(err).Any("name", name).Msg("failed to insert tenant")
		return "", err
	}
	return id, nil
}
