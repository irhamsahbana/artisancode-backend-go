package seeds

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Seed) jobPositionSeed() {
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

	type jobPositionDef struct {
		Name  string
		Grade string
	}

	positions := []jobPositionDef{
		{Name: "Software Engineer", Grade: "SE-1"},
		{Name: "Senior Software Engineer", Grade: "SE-2"},
		{Name: "Tech Lead", Grade: "TL-1"},
		{Name: "HR Manager", Grade: "HR-1"},
		{Name: "Finance Staff", Grade: "FN-1"},
		{Name: "UI/UX Designer", Grade: "DX-1"},
		{Name: "Project Manager", Grade: "PM-1"},
	}

	for _, jp := range positions {
		if _, err := upsertJobPosition(tx, tenantID, jp.Name, jp.Grade); err != nil {
			log.Error().Err(err).Str("name", jp.Name).Msg("Error upserting job position")
			return
		}
	}

	log.Info().Msg("job_positions seeded successfully!")
}

func upsertJobPosition(tx *sqlx.Tx, tenantID string, name string, grade string) (string, error) {
	var id string
	selectQuery := `
		SELECT id
		FROM job_positions
		WHERE name = ? AND tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT 1
	`
	if err := tx.Get(&id, tx.Rebind(selectQuery), name, tenantID); err == nil {
		return id, nil
	}

	insertQuery := `
		INSERT INTO job_positions (tenant_id, name, grade)
		VALUES (?, ?, ?)
		RETURNING id
	`
	if err := tx.Get(&id, tx.Rebind(insertQuery), tenantID, name, grade); err != nil {
		return "", err
	}
	return id, nil
}