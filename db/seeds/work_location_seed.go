package seeds

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Seed) workLocationSeed() {
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

	type workLocationDef struct {
		Name         string
		OrgUnitName  string
		Address      *string
		Timezone     string
		Latitude     *float64
		Longitude    *float64
		RadiusMeters *int
	}

	hqLat := -6.2088
	hqLon := 106.8456
	hqRadius := 200
	cabLat := -6.9175
	cabLon := 107.6191
	cabRadius := 150
	hqAddr := "Jl. Sudirman No. 1, Jakarta"
	cabAddr := "Jl. Asia Afrika No. 10, Bandung"

	locations := []workLocationDef{
		{
			Name: "Kantor Pusat Jakarta", OrgUnitName: "Company", Address: &hqAddr,
			Timezone: "Asia/Jakarta", Latitude: &hqLat, Longitude: &hqLon, RadiusMeters: &hqRadius,
		},
		{
			Name: "Kantor Cabang Bandung", OrgUnitName: "Division A", Address: &cabAddr,
			Timezone: "Asia/Jakarta", Latitude: &cabLat, Longitude: &cabLon, RadiusMeters: &cabRadius,
		},
		{
			Name: "Remote Area", OrgUnitName: "", Address: nil,
			Timezone: "Asia/Jakarta", Latitude: nil, Longitude: nil, RadiusMeters: nil,
		},
	}

	for _, wl := range locations {
		var orgUnitID *string
		if wl.OrgUnitName != "" {
			id, err := getOrgUnitIDByName(tx, tenantID, wl.OrgUnitName)
			if err != nil {
				log.Error().Err(err).Str("org_unit", wl.OrgUnitName).Msg("Error finding org unit")
				return
			}
			orgUnitID = &id
		}

		if _, err := upsertWorkLocation(tx, tenantID, orgUnitID, wl.Name, wl.Address, wl.Timezone, wl.Latitude, wl.Longitude, wl.RadiusMeters); err != nil {
			log.Error().Err(err).Str("name", wl.Name).Msg("Error upserting work location")
			return
		}
	}

	log.Info().Msg("work_locations seeded successfully!")
}

func getOrgUnitIDByName(tx *sqlx.Tx, _ string, name string) (string, error) {
	var id string
	var defaultTenantID *string
	query := `
		SELECT id
		FROM org_units
		WHERE name = ? AND tenant_id IS NOT DISTINCT FROM ? AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT 1
	`
	if err := tx.Get(&id, tx.Rebind(query), name, defaultTenantID); err != nil {
		return "", err
	}
	return id, nil
}

func upsertWorkLocation(tx *sqlx.Tx, tenantID string, orgUnitID *string, name string, address *string, timezone string, latitude *float64, longitude *float64, radiusMeters *int) (string, error) {
	var id string
	selectQuery := `
		SELECT id
		FROM work_locations
		WHERE name = ? AND tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT 1
	`
	if err := tx.Get(&id, tx.Rebind(selectQuery), name, tenantID); err == nil {
		return id, nil
	}

	insertQuery := `
		INSERT INTO work_locations (tenant_id, org_unit_id, name, address, timezone, latitude, longitude, radius_meters)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	if err := tx.Get(&id, tx.Rebind(insertQuery), tenantID, orgUnitID, name, address, timezone, latitude, longitude, radiusMeters); err != nil {
		return "", err
	}
	return id, nil
}