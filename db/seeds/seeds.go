package seeds

import (
	"context"
	"flag"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const (
	seedTableTenants                    = "tenants"
	seedTablePermissions                = "permissions"
	seedTableRoles                      = "roles"
	seedTableRolePerms                  = "role_permissions"
	seedTableInternalTemplateRoles      = "internal_template_roles"
	seedTableInternalTemplatePerms      = "internal_template_permissions"
	seedTableInternalTemplateRolePerms  = "internal_template_role_permissions"
	seedTableOrgUnits                   = "org_units"
	seedTableUsers                      = "users"
	seedTableUserRoles                  = "user_roles"
	seedTableJobPositions               = "job_positions"
	seedTableWorkLocations              = "work_locations"
	seedTableWorkShifts                 = "work_shifts"
	seedTableEmployees                  = "employees"
	seedTableInternalUsers              = "internal_users"
	seedTableInternalCurrencies         = "internal_currencies"
	seedTableInternalProviderCurrencies = "internal_payment_provider_currencies"
	seedTableInternalProducts           = "internal_products"
	seedTableInternalProductPricings    = "internal_product_pricings"
	seedTableInternalProductPrices      = "internal_product_prices"
)

var seedGroups = map[string][]string{
	"tenant":            {seedTableTenants},
	"company":           {seedTableTenants},
	"rbac":              {seedTableInternalTemplatePerms, seedTableInternalTemplateRoles, seedTableInternalTemplateRolePerms, seedTablePermissions, seedTableRoles, seedTableRolePerms},
	"org_unit":          {seedTableOrgUnits},
	"owner_user":        {seedTableTenants, seedTableRoles, seedTableUsers, seedTableUserRoles},
	"job_position":      {seedTableTenants, seedTableJobPositions},
	"work_location":     {seedTableOrgUnits, seedTableTenants, seedTableWorkLocations},
	"work_shift":        {seedTableTenants, seedTableWorkShifts},
	"employee":          {seedTableTenants, seedTablePermissions, seedTableRoles, seedTableRolePerms, seedTableOrgUnits, seedTableUsers, seedTableUserRoles, seedTableJobPositions, seedTableWorkLocations, seedTableWorkShifts, seedTableEmployees},
	"template":          {seedTableInternalTemplatePerms, seedTableInternalTemplateRoles, seedTableInternalTemplateRolePerms, seedTableOrgUnits},
	"all":               {seedTableTenants, seedTableInternalTemplatePerms, seedTableInternalTemplateRoles, seedTableInternalTemplateRolePerms, seedTablePermissions, seedTableRoles, seedTableRolePerms, seedTableUsers, seedTableUserRoles, seedTableInternalUsers, seedTableInternalCurrencies, seedTableInternalProviderCurrencies, seedTableInternalProducts, seedTableInternalProductPricings, seedTableInternalProductPrices},
	"all_with_template": {seedTableTenants, seedTableInternalTemplatePerms, seedTableInternalTemplateRoles, seedTableInternalTemplateRolePerms, seedTablePermissions, seedTableRoles, seedTableRolePerms, seedTableOrgUnits, seedTableUsers, seedTableUserRoles, seedTableJobPositions, seedTableWorkLocations, seedTableWorkShifts, seedTableEmployees, seedTableInternalUsers, seedTableInternalCurrencies, seedTableInternalProviderCurrencies, seedTableInternalProducts, seedTableInternalProductPricings, seedTableInternalProductPrices},
	"internal_template": {seedTableInternalTemplatePerms, seedTableInternalTemplateRoles, seedTableInternalTemplateRolePerms},
	"internal_users":    {seedTableInternalUsers},
	"internal_currency": {seedTableInternalCurrencies},
	"internal_catalog":  {seedTableInternalCurrencies, seedTableInternalProviderCurrencies, seedTableInternalProducts, seedTableInternalProductPricings, seedTableInternalProductPrices},
	"all_internal":      {seedTableInternalTemplatePerms, seedTableInternalTemplateRoles, seedTableInternalTemplateRolePerms, seedTableInternalUsers, seedTableInternalCurrencies, seedTableInternalProviderCurrencies, seedTableInternalProducts, seedTableInternalProductPricings, seedTableInternalProductPrices},
}

type Seed struct {
	db *sqlx.DB
}

func newSeed(db *sqlx.DB) Seed {
	return Seed{db: db}
}

func Execute(db *sqlx.DB, table string, total int) error {
	_ = total

	seed := newSeed(db)
	return seed.run(table)
}

func (s *Seed) run(table string) error {
	order, ok := seedGroups[table]
	if !ok {
		log.Warn().Str("table", table).Msg("No seed to run")
		return nil
	}

	state, err := newCSVSeedState()
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}

	for _, tableName := range order {
		log.Info().Str("table", tableName).Msg("Seeding table")
		err = s.seedTable(tx, state, tableName)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("Failed to rollback seed transaction")
			}
			return err
		}
		log.Info().Str("table", tableName).Msg("Seeded table successfully")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}

	if err := state.writeModifiedFiles(); err != nil {
		return err
	}

	return nil
}

func (s *Seed) seedTable(tx *sqlx.Tx, state *csvSeedState, tableName string) error {
	switch tableName {
	case seedTableTenants:
		return s.seedTenants(tx, state)
	case seedTablePermissions:
		return s.seedPermissions(tx, state)
	case seedTableRoles:
		return s.seedRoles(tx, state)
	case seedTableRolePerms:
		return s.seedRolePermissions(tx, state)
	case seedTableInternalTemplateRoles:
		return s.seedInternalTemplateRoles(tx, state)
	case seedTableInternalTemplatePerms:
		return s.seedInternalTemplatePermissions(tx, state)
	case seedTableInternalTemplateRolePerms:
		return s.seedInternalTemplateRolePermissions(tx, state)
	case seedTableOrgUnits:
		return s.seedOrgUnits(tx, state)
	case seedTableUsers:
		return s.seedUsers(tx, state)
	case seedTableUserRoles:
		return s.seedUserRoles(tx, state)
	case seedTableJobPositions:
		return s.seedJobPositions(tx, state)
	case seedTableWorkLocations:
		return s.seedWorkLocations(tx, state)
	case seedTableWorkShifts:
		return s.seedWorkShifts(tx, state)
	case seedTableEmployees:
		return s.seedEmployees(tx, state)
	case seedTableInternalUsers:
		return s.seedInternalUsers(tx, state)
	case seedTableInternalCurrencies:
		return s.seedInternalCurrencies(tx, state)
	case seedTableInternalProviderCurrencies:
		return s.seedInternalPaymentProviderCurrencies(tx, state)
	case seedTableInternalProducts:
		return s.seedInternalProducts(tx, state)
	case seedTableInternalProductPricings:
		return s.seedInternalProductPricings(tx, state)
	case seedTableInternalProductPrices:
		return s.seedInternalProductPrices(tx, state)
	default:
		return fmt.Errorf("unsupported seed table %q", tableName)
	}
}

func BindSeedFlags(cmd *flag.FlagSet) (*string, *int) {
	table := cmd.String("table", "", "seed to run")
	total := cmd.Int("total", 1, "total of records to seed")
	return table, total
}
