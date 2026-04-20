package seeds

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type upsertConfig struct {
	table              string
	columns            []string
	matchColumns       []string
	hasUpdatedAt       bool
	supportsSoftDelete bool
}

func upsertRecord(tx *sqlx.Tx, cfg upsertConfig, rowID string, values map[string]any) (string, bool, error) {
	normalizedID := strings.TrimSpace(rowID)
	if normalizedID != "" {
		if err := validateUUIDv7(normalizedID); err != nil {
			return "", false, err
		}
		exists, err := recordExistsByID(tx, cfg.table, normalizedID)
		if err != nil {
			return "", false, err
		}
		if exists {
			if err := updateRecord(tx, cfg, normalizedID, values); err != nil {
				return "", false, err
			}
			return normalizedID, false, nil
		}
		if err := insertRecord(tx, cfg, normalizedID, values); err != nil {
			return "", false, err
		}
		return normalizedID, false, nil
	}

	existingID, err := findExistingID(tx, cfg, values)
	if err != nil {
		return "", false, err
	}
	if existingID != "" {
		if err := updateRecord(tx, cfg, existingID, values); err != nil {
			return "", false, err
		}
		return existingID, true, nil
	}

	generatedID, err := newUUIDv7()
	if err != nil {
		return "", false, err
	}
	if err := insertRecord(tx, cfg, generatedID, values); err != nil {
		return "", false, err
	}
	return generatedID, true, nil
}

func recordExistsByID(tx *sqlx.Tx, table string, id string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = ?)", table)
	if err := tx.Get(&exists, tx.Rebind(query), id); err != nil {
		return false, fmt.Errorf("check existing id for %s: %w", table, err)
	}
	return exists, nil
}

func findExistingID(tx *sqlx.Tx, cfg upsertConfig, values map[string]any) (string, error) {
	clauses := make([]string, 0, len(cfg.matchColumns)+1)
	args := make([]any, 0, len(cfg.matchColumns))
	for _, column := range cfg.matchColumns {
		clauses = append(clauses, column+" IS NOT DISTINCT FROM ?")
		args = append(args, values[column])
	}
	if cfg.supportsSoftDelete {
		clauses = append(clauses, "deleted_at IS NULL")
	}

	query := fmt.Sprintf("SELECT id FROM %s WHERE %s ORDER BY id ASC LIMIT 1", cfg.table, strings.Join(clauses, " AND "))
	var id string
	err := tx.Get(&id, tx.Rebind(query), args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("find existing %s row: %w", cfg.table, err)
	}
	return id, nil
}

func insertRecord(tx *sqlx.Tx, cfg upsertConfig, id string, values map[string]any) error {
	columns := append([]string{"id"}, cfg.columns...)
	placeholders := make([]string, 0, len(columns))
	args := make([]any, 0, len(columns))
	args = append(args, id)
	placeholders = append(placeholders, "?")

	for _, column := range cfg.columns {
		placeholders = append(placeholders, "?")
		args = append(args, values[column])
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		cfg.table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	if _, err := tx.Exec(tx.Rebind(query), args...); err != nil {
		return fmt.Errorf("insert %s row: %w", cfg.table, err)
	}
	return nil
}

func updateRecord(tx *sqlx.Tx, cfg upsertConfig, id string, values map[string]any) error {
	setClauses := make([]string, 0, len(cfg.columns)+2)
	args := make([]any, 0, len(cfg.columns)+1)

	for _, column := range cfg.columns {
		setClauses = append(setClauses, column+" = ?")
		args = append(args, values[column])
	}
	if cfg.hasUpdatedAt {
		setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")
	}
	if cfg.supportsSoftDelete {
		setClauses = append(setClauses, "deleted_at = NULL")
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", cfg.table, strings.Join(setClauses, ", "))
	if _, err := tx.Exec(tx.Rebind(query), args...); err != nil {
		return fmt.Errorf("update %s row: %w", cfg.table, err)
	}
	return nil
}

func newUUIDv7() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("generate uuidv7: %w", err)
	}
	return id.String(), nil
}

func validateUUIDv7(value string) error {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return fmt.Errorf("invalid uuid %q: %w", value, err)
	}
	if parsed.Version() != 7 {
		return fmt.Errorf("uuid %q is not version 7", value)
	}
	return nil
}

func requiredCSVValue(row map[string]string, column string) (string, error) {
	value := strings.TrimSpace(row[column])
	if value == "" {
		return "", fmt.Errorf("column %q is required", column)
	}
	return value, nil
}

func nullableCSVValue(row map[string]string, column string) *string {
	value := strings.TrimSpace(row[column])
	if value == "" {
		return nil
	}
	return &value
}

func optionalIntValue(row map[string]string, column string) (*int, error) {
	raw := strings.TrimSpace(row[column])
	if raw == "" {
		return nil, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return nil, fmt.Errorf("parse %s as int: %w", column, err)
	}
	return &value, nil
}

func requiredIntValue(row map[string]string, column string) (int, error) {
	raw, err := requiredCSVValue(row, column)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("parse %s as int: %w", column, err)
	}
	return value, nil
}

func optionalFloatValue(row map[string]string, column string) (*float64, error) {
	raw := strings.TrimSpace(row[column])
	if raw == "" {
		return nil, nil
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return nil, fmt.Errorf("parse %s as float: %w", column, err)
	}
	return &value, nil
}

func optionalTimeValue(row map[string]string, column string) (*time.Time, error) {
	raw := strings.TrimSpace(row[column])
	if raw == "" {
		return nil, nil
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		value, err := time.Parse(layout, raw)
		if err == nil {
			return &value, nil
		}
	}

	return nil, fmt.Errorf("parse %s as time: unsupported format %q", column, raw)
}
