package seeds

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Seed) workShiftSeed() {
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

	type workShiftDef struct {
		Name             string
		Timezone         string
		StartTime        string
		EndTime          string
		GracePeriodMins  int
	}

	shifts := []workShiftDef{
		{Name: "Shift Pagi (08:00-17:00)", Timezone: "Asia/Jakarta", StartTime: "08:00", EndTime: "17:00", GracePeriodMins: 15},
		{Name: "Shift Siang (13:00-22:00)", Timezone: "Asia/Jakarta", StartTime: "13:00", EndTime: "22:00", GracePeriodMins: 15},
		{Name: "Shift Malam (22:00-06:00)", Timezone: "Asia/Jakarta", StartTime: "22:00", EndTime: "06:00", GracePeriodMins: 15},
		{Name: "Flexible", Timezone: "Asia/Jakarta", StartTime: "09:00", EndTime: "18:00", GracePeriodMins: 30},
	}

	for _, ws := range shifts {
		if _, err := upsertWorkShift(tx, tenantID, ws.Name, ws.Timezone, ws.StartTime, ws.EndTime, ws.GracePeriodMins); err != nil {
			log.Error().Err(err).Str("name", ws.Name).Msg("Error upserting work shift")
			return
		}
	}

	log.Info().Msg("work_shifts seeded successfully!")
}

func upsertWorkShift(tx *sqlx.Tx, tenantID string, name string, timezone string, startTime string, endTime string, gracePeriodMins int) (string, error) {
	var id string
	selectQuery := `
		SELECT id
		FROM work_shifts
		WHERE name = ? AND tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT 1
	`
	if err := tx.Get(&id, tx.Rebind(selectQuery), name, tenantID); err == nil {
		return id, nil
	}

	insertQuery := `
		INSERT INTO work_shifts (tenant_id, name, timezone, start_time, end_time, grace_period_minutes)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	if err := tx.Get(&id, tx.Rebind(insertQuery), tenantID, name, timezone, startTime, endTime, gracePeriodMins); err != nil {
		return "", err
	}
	return id, nil
}