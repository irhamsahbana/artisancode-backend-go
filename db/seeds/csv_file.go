package seeds

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type csvSeedFile struct {
	logicalName string
	path        string
	header      []string
	rows        []map[string]string
	modified    bool
}

type csvSeedState struct {
	files map[string]*csvSeedFile
	ids   map[string]map[string]string
}

var csvSeedFiles = map[string]string{
	seedTableTenants:                 "tenants.csv",
	seedTablePermissions:             "permissions.csv",
	seedTableRoles:                   "roles.csv",
	seedTableRolePerms:               "role_permissions.csv",
	seedTableOrgUnits:                "org_units.csv",
	seedTableUsers:                   "users.csv",
	seedTableUserRoles:               "user_roles.csv",
	seedTableJobPositions:            "job_positions.csv",
	seedTableWorkLocations:           "work_locations.csv",
	seedTableWorkShifts:              "work_shifts.csv",
	seedTableEmployees:               "employees.csv",
	seedTableInternalUsers:           "internal_users.csv",
	seedTableInternalProducts:        "internal_products.csv",
	seedTableInternalProductPricings: "internal_product_pricings.csv",
	seedTableInternalProductPrices:   "internal_product_prices.csv",
}

const helperColumnPrefix = "helper__"

var helperColumnsByFile = map[string]map[string]struct{}{
	seedTablePermissions: {
		"tenant_code": {},
	},
	seedTableRoles: {
		"tenant_code": {},
	},
	seedTableRolePerms: {
		"tenant_code":     {},
		"role_name":       {},
		"permission_name": {},
	},
	seedTableOrgUnits: {
		"tenant_code": {},
		"parent_code": {},
	},
	seedTableUsers: {
		"tenant_code":         {},
		"company_tenant_code": {},
		"company_code":        {},
		"password_plain":      {},
	},
	seedTableUserRoles: {
		"tenant_code": {},
		"user_email":  {},
		"role_name":   {},
	},
	seedTableJobPositions: {
		"tenant_code": {},
	},
	seedTableWorkLocations: {
		"tenant_code":          {},
		"org_unit_tenant_code": {},
		"org_unit_code":        {},
	},
	seedTableWorkShifts: {
		"tenant_code": {},
	},
	seedTableEmployees: {
		"tenant_code":              {},
		"user_email":               {},
		"org_unit_tenant_code":     {},
		"org_unit_code":            {},
		"job_position_tenant_code": {},
		"job_position_name":        {},
		"location_tenant_code":     {},
		"location_name":            {},
		"shift_tenant_code":        {},
		"shift_name":               {},
	},
	seedTableInternalUsers: {
		"password_plain": {},
	},
	seedTableInternalProductPricings: {
		"product_code": {},
	},
	seedTableInternalProductPrices: {
		"product_code": {},
		"pricing_code": {},
	},
}

func newCSVSeedState() (*csvSeedState, error) {
	baseDir, err := csvSeedDataDir()
	if err != nil {
		return nil, err
	}

	state := &csvSeedState{
		files: make(map[string]*csvSeedFile, len(csvSeedFiles)),
		ids:   make(map[string]map[string]string),
	}

	for logicalName, fileName := range csvSeedFiles {
		filePath := filepath.Join(baseDir, fileName)
		seedFile, err := loadCSVSeedFile(logicalName, filePath)
		if err != nil {
			return nil, err
		}
		state.files[logicalName] = seedFile
	}

	state.rebuildLookups()

	return state, nil
}

func csvSeedDataDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("resolve seed data directory")
	}

	return filepath.Join(filepath.Dir(filename), "data"), nil
}

func loadCSVSeedFile(logicalName string, path string) (*csvSeedFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open seed csv %s: %w", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read seed csv %s: %w", path, err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("seed csv %s is empty", path)
	}

	header := make([]string, len(records[0]))
	copy(header, records[0])
	if err := validateCSVHeader(logicalName, header); err != nil {
		return nil, err
	}

	rows := make([]map[string]string, 0, len(records)-1)
	for _, record := range records[1:] {
		row := make(map[string]string, len(header))
		isEmpty := true
		for idx, column := range header {
			logicalColumn := normalizeCSVColumnName(column)
			value := ""
			if idx < len(record) {
				value = strings.TrimSpace(record[idx])
			}
			if value != "" {
				isEmpty = false
			}
			row[column] = value
			if logicalColumn != column {
				row[logicalColumn] = value
			}
		}
		if isEmpty {
			continue
		}
		rows = append(rows, row)
	}

	return &csvSeedFile{
		logicalName: logicalName,
		path:        path,
		header:      header,
		rows:        rows,
	}, nil
}

func (s *csvSeedState) rebuildLookups() {
	s.ids = make(map[string]map[string]string)

	register := func(tableName, key, id string) {
		key = strings.TrimSpace(key)
		id = strings.TrimSpace(id)
		if key == "" || id == "" {
			return
		}
		if _, ok := s.ids[tableName]; !ok {
			s.ids[tableName] = map[string]string{}
		}
		s.ids[tableName][key] = id
	}

	for _, row := range s.files[seedTableTenants].rows {
		register(seedTableTenants, tenantLookupKey(row["code"]), row["id"])
	}
	for _, row := range s.files[seedTablePermissions].rows {
		register(seedTablePermissions, scopedLookupKey(row["tenant_code"], row["name"]), row["id"])
	}
	for _, row := range s.files[seedTableRoles].rows {
		register(seedTableRoles, scopedLookupKey(row["tenant_code"], row["name"]), row["id"])
	}
	for _, row := range s.files[seedTableOrgUnits].rows {
		register(seedTableOrgUnits, scopedLookupKey(row["tenant_code"], row["code"]), row["id"])
	}
	for _, row := range s.files[seedTableUsers].rows {
		register(seedTableUsers, scopedLookupKey(row["tenant_code"], row["email"]), row["id"])
	}
	for _, row := range s.files[seedTableJobPositions].rows {
		register(seedTableJobPositions, scopedLookupKey(row["tenant_code"], row["name"]), row["id"])
	}
	for _, row := range s.files[seedTableWorkLocations].rows {
		register(seedTableWorkLocations, scopedLookupKey(row["tenant_code"], row["name"]), row["id"])
	}
	for _, row := range s.files[seedTableWorkShifts].rows {
		register(seedTableWorkShifts, scopedLookupKey(row["tenant_code"], row["name"]), row["id"])
	}
	for _, row := range s.files[seedTableEmployees].rows {
		register(seedTableEmployees, scopedLookupKey(row["tenant_code"], row["employee_no"]), row["id"])
	}
	for _, row := range s.files[seedTableInternalUsers].rows {
		register(seedTableInternalUsers, tenantLookupKey(row["email"]), row["id"])
	}
	for _, row := range s.files[seedTableInternalProducts].rows {
		register(seedTableInternalProducts, row["code"], row["id"])
	}
	for _, row := range s.files[seedTableInternalProductPricings].rows {
		register(seedTableInternalProductPricings, row["product_code"]+"::"+row["code"], row["id"])
	}
	for _, row := range s.files[seedTableInternalProductPrices].rows {
		register(seedTableInternalProductPrices, row["product_code"]+"::"+row["pricing_code"]+"::"+row["currency_code"]+"::"+row["started_at"], row["id"])
	}
}

func (s *csvSeedState) lookupID(tableName string, key string) (string, bool) {
	tableEntries, ok := s.ids[tableName]
	if !ok {
		return "", false
	}
	id, ok := tableEntries[key]
	return id, ok
}

func (s *csvSeedState) registerID(tableName string, key string, id string) {
	if _, ok := s.ids[tableName]; !ok {
		s.ids[tableName] = map[string]string{}
	}
	s.ids[tableName][key] = id
}

func (s *csvSeedState) writeModifiedFiles() error {
	for _, seedFile := range s.files {
		if !seedFile.modified {
			continue
		}
		if err := seedFile.write(); err != nil {
			return err
		}
	}
	return nil
}

func (f *csvSeedFile) write() error {
	file, err := os.Create(f.path)
	if err != nil {
		return fmt.Errorf("write seed csv %s: %w", f.path, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write(f.header); err != nil {
		return fmt.Errorf("write csv header %s: %w", f.path, err)
	}

	for _, row := range f.rows {
		record := make([]string, 0, len(f.header))
		for _, column := range f.header {
			value := row[column]
			if value == "" {
				value = row[normalizeCSVColumnName(column)]
			}
			record = append(record, value)
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("write csv row %s: %w", f.path, err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush csv %s: %w", f.path, err)
	}

	return nil
}

func tenantLookupKey(code string) string {
	return strings.TrimSpace(code)
}

func scopedLookupKey(scope string, value string) string {
	scope = strings.TrimSpace(scope)
	if scope == "" {
		scope = "__global__"
	}
	return scope + "::" + strings.TrimSpace(value)
}

func normalizeCSVColumnName(column string) string {
	column = strings.TrimSpace(column)
	return strings.TrimPrefix(column, helperColumnPrefix)
}

func validateCSVHeader(logicalName string, header []string) error {
	helperColumns, ok := helperColumnsByFile[logicalName]
	if !ok {
		return nil
	}

	for _, column := range header {
		column = strings.TrimSpace(column)
		normalizedColumn := normalizeCSVColumnName(column)
		_, isHelper := helperColumns[normalizedColumn]
		if isHelper && !strings.HasPrefix(column, helperColumnPrefix) {
			return fmt.Errorf("seed csv %s: helper column %q must use %q prefix", logicalName, normalizedColumn, helperColumnPrefix)
		}
	}

	return nil
}
