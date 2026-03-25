package seeds

import (
	"context"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Seed) orgUnitSeed() {
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

	var defaultTenantID *string

	type orgUnitDef struct {
		Name     string
		Parent   string
		Category string
		Code     string
	}

	orgUnits := []orgUnitDef{
		{Name: "Company", Parent: "", Category: "company", Code: "COMPANY"},
		{Name: "Division A", Parent: "Company", Category: "division", Code: "DIV-A"},
		{Name: "Division B", Parent: "Company", Category: "division", Code: "DIV-B"},
		{Name: "Department HR", Parent: "Division A", Category: "department", Code: "DEPT-HR"},
		{Name: "Department Finance", Parent: "Division A", Category: "department", Code: "DEPT-FIN"},
		{Name: "Department IT", Parent: "Division B", Category: "department", Code: "DEPT-IT"},
		{Name: "Unit Programming", Parent: "Department IT", Category: "unit", Code: "UNIT-PROG"},
		{Name: "Unit Design", Parent: "Department IT", Category: "unit", Code: "UNIT-DES"},
	}

	idMapping := make(map[string]string)

	for _, ou := range orgUnits {
		id, err := upsertOrgUnit(tx, defaultTenantID, ou.Name, ou.Code, ou.Category)
		if err != nil {
			log.Error().Err(err).Str("name", ou.Name).Msg("Error upserting org unit")
			return
		}
		idMapping[ou.Name] = id
	}

	for _, ou := range orgUnits {
		if ou.Parent == "" {
			continue
		}
		parentID := idMapping[ou.Parent]
		orgUnitID := idMapping[ou.Name]
		if parentID == "" || orgUnitID == "" {
			continue
		}
		if err := updateOrgUnitParent(tx, orgUnitID, parentID); err != nil {
			log.Error().Err(err).Str("name", ou.Name).Msg("Error updating org unit parent")
			return
		}
	}

	log.Info().Msg("org_unit template seeded successfully!")
}

func makeCode(name string) string {
	var result []byte
	for i, r := range name {
		if r == ' ' || r == '-' {
			continue
		}
		if i < 4 {
			result = append(result, byte(r))
		}
	}
	if len(result) > 4 {
		result = result[:4]
	}
	return string(result)
}

func ensureMaxLen(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxLen])
}

func upsertOrgUnit(tx *sqlx.Tx, tenantID *string, name string, code string, orgCategory string) (string, error) {
	if code == "" {
		code = ensureMaxLen(makeCode(name), 8)
	}

	var id string
	selectQuery := `
		SELECT id
		FROM org_units
		WHERE name = ? AND tenant_id IS NOT DISTINCT FROM ? AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT 1
	`
	if err := tx.Get(&id, tx.Rebind(selectQuery), name, tenantID); err == nil {
		return id, nil
	}

	query := `
		INSERT INTO org_units (tenant_id, name, code, category, config)
		VALUES (?, ?, ?, ?, '{}')
		RETURNING id
	`
	if err := tx.Get(&id, tx.Rebind(query), tenantID, name, code, orgCategory); err != nil {
		return "", err
	}
	return id, nil
}

func updateOrgUnitParent(tx *sqlx.Tx, orgUnitID string, parentID string) error {
	query := `
		UPDATE org_units
		SET parent_id = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := tx.ExecContext(context.Background(), tx.Rebind(query), parentID, orgUnitID)
	return err
}
