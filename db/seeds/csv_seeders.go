package seeds

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func (s *Seed) seedTenants(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableTenants]
	cfg := upsertConfig{
		table:              seedTableTenants,
		columns:            []string{"name", "code", "config"},
		matchColumns:       []string{"code"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		code, err := requiredCSVValue(row, "code")
		if err != nil {
			return fmt.Errorf("seed tenants: %w", err)
		}
		name, err := requiredCSVValue(row, "name")
		if err != nil {
			return fmt.Errorf("seed tenants %s: %w", code, err)
		}

		values := map[string]any{
			"name":   name,
			"code":   code,
			"config": defaultJSON(row["config"]),
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed tenants %s: %w", code, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableTenants, tenantLookupKey(code), id)
	}

	return nil
}

func (s *Seed) seedPermissions(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTablePermissions]
	cfg := upsertConfig{
		table:              seedTablePermissions,
		columns:            []string{"tenant_id", "name", "description"},
		matchColumns:       []string{"tenant_id", "name"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		name, err := requiredCSVValue(row, "name")
		if err != nil {
			return fmt.Errorf("seed permissions: %w", err)
		}

		tenantID, err := resolveTenantID(state, row["tenant_code"])
		if err != nil {
			return fmt.Errorf("seed permissions %s: %w", name, err)
		}

		values := map[string]any{
			"tenant_id":   tenantID,
			"name":        name,
			"description": nullableStringValue(row["description"]),
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed permissions %s: %w", name, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTablePermissions, scopedLookupKey(row["tenant_code"], name), id)
	}

	return nil
}

func (s *Seed) seedRoles(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableRoles]
	cfg := upsertConfig{
		table:              seedTableRoles,
		columns:            []string{"tenant_id", "name"},
		matchColumns:       []string{"tenant_id", "name"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		name, err := requiredCSVValue(row, "name")
		if err != nil {
			return fmt.Errorf("seed roles: %w", err)
		}

		tenantID, err := resolveTenantID(state, row["tenant_code"])
		if err != nil {
			return fmt.Errorf("seed roles %s: %w", name, err)
		}

		values := map[string]any{
			"tenant_id": tenantID,
			"name":      name,
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed roles %s: %w", name, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableRoles, scopedLookupKey(row["tenant_code"], name), id)
	}

	return nil
}

func (s *Seed) seedRolePermissions(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableRolePerms]

	for _, row := range file.rows {
		roleName, err := requiredCSVValue(row, "role_name")
		if err != nil {
			return fmt.Errorf("seed role_permissions: %w", err)
		}
		permissionName, err := requiredCSVValue(row, "permission_name")
		if err != nil {
			return fmt.Errorf("seed role_permissions %s/%s: %w", roleName, permissionName, err)
		}

		roleID, err := resolveScopedID(state, seedTableRoles, row["tenant_code"], roleName, "role")
		if err != nil {
			return fmt.Errorf("seed role_permissions %s/%s: %w", roleName, permissionName, err)
		}
		permissionID, err := resolveScopedID(state, seedTablePermissions, row["tenant_code"], permissionName, "permission")
		if err != nil {
			return fmt.Errorf("seed role_permissions %s/%s: %w", roleName, permissionName, err)
		}
		tenantID, err := resolveTenantID(state, row["tenant_code"])
		if err != nil {
			return fmt.Errorf("seed role_permissions %s/%s: %w", roleName, permissionName, err)
		}

		if strings.TrimSpace(row["id"]) == "" {
			id, err := newUUIDv7()
			if err != nil {
				return fmt.Errorf("seed role_permissions %s/%s: %w", roleName, permissionName, err)
			}
			row["id"] = id
			file.modified = true
		} else if err := validateUUIDv7(row["id"]); err != nil {
			return fmt.Errorf("seed role_permissions %s/%s: %w", roleName, permissionName, err)
		}

		query := `
			INSERT INTO role_permissions (role_id, permission_id, tenant_id)
			VALUES (?, ?, ?)
			ON CONFLICT (role_id, permission_id) DO UPDATE SET
				tenant_id = EXCLUDED.tenant_id
		`
		if _, err := tx.Exec(tx.Rebind(query), roleID, permissionID, tenantID); err != nil {
			return fmt.Errorf("seed role_permissions %s/%s: %w", roleName, permissionName, err)
		}
	}

	return nil
}

func (s *Seed) seedInternalTemplatePermissions(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableInternalTemplatePerms]
	cfg := upsertConfig{
		table:              seedTableInternalTemplatePerms,
		columns:            []string{"name", "description"},
		matchColumns:       []string{"name"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		name, err := requiredCSVValue(row, "name")
		if err != nil {
			return fmt.Errorf("seed internal_template_permissions: %w", err)
		}

		values := map[string]any{
			"name":        name,
			"description": nullableStringValue(row["description"]),
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed internal_template_permissions %s: %w", name, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableInternalTemplatePerms, name, id)
	}

	return nil
}

func (s *Seed) seedInternalTemplateRoles(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableInternalTemplateRoles]
	cfg := upsertConfig{
		table:              seedTableInternalTemplateRoles,
		columns:            []string{"name"},
		matchColumns:       []string{"name"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		name, err := requiredCSVValue(row, "name")
		if err != nil {
			return fmt.Errorf("seed internal_template_roles: %w", err)
		}

		values := map[string]any{
			"name": name,
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed internal_template_roles %s: %w", name, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableInternalTemplateRoles, name, id)
	}

	return nil
}

func (s *Seed) seedInternalTemplateRolePermissions(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableInternalTemplateRolePerms]

	for _, row := range file.rows {
		roleName, err := requiredCSVValue(row, "role_name")
		if err != nil {
			return fmt.Errorf("seed internal_template_role_permissions: %w", err)
		}
		permissionName, err := requiredCSVValue(row, "permission_name")
		if err != nil {
			return fmt.Errorf("seed internal_template_role_permissions %s/%s: %w", roleName, permissionName, err)
		}

		roleID, err := resolveTemplateID(state, seedTableInternalTemplateRoles, roleName, "internal template role")
		if err != nil {
			return fmt.Errorf("seed internal_template_role_permissions %s/%s: %w", roleName, permissionName, err)
		}
		permissionID, err := resolveTemplateID(state, seedTableInternalTemplatePerms, permissionName, "internal template permission")
		if err != nil {
			return fmt.Errorf("seed internal_template_role_permissions %s/%s: %w", roleName, permissionName, err)
		}

		query := `
			INSERT INTO internal_template_role_permissions (role_id, permission_id)
			VALUES (?, ?)
			ON CONFLICT (role_id, permission_id) DO NOTHING
		`
		if _, err := tx.Exec(tx.Rebind(query), roleID, permissionID); err != nil {
			return fmt.Errorf("seed internal_template_role_permissions %s/%s: %w", roleName, permissionName, err)
		}
	}

	return nil
}

func (s *Seed) seedOrgUnits(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableOrgUnits]
	cfg := upsertConfig{
		table:              seedTableOrgUnits,
		columns:            []string{"tenant_id", "parent_id", "name", "code", "category", "config"},
		matchColumns:       []string{"tenant_id", "code"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		code, err := requiredCSVValue(row, "code")
		if err != nil {
			return fmt.Errorf("seed org_units: %w", err)
		}
		name, err := requiredCSVValue(row, "name")
		if err != nil {
			return fmt.Errorf("seed org_units %s: %w", code, err)
		}
		category, err := requiredCSVValue(row, "category")
		if err != nil {
			return fmt.Errorf("seed org_units %s: %w", code, err)
		}

		tenantID, err := resolveTenantID(state, row["tenant_code"])
		if err != nil {
			return fmt.Errorf("seed org_units %s: %w", code, err)
		}
		parentID, err := resolveScopedNullableID(state, seedTableOrgUnits, row["tenant_code"], row["parent_code"], "parent org unit")
		if err != nil {
			return fmt.Errorf("seed org_units %s: %w", code, err)
		}

		values := map[string]any{
			"tenant_id": tenantID,
			"parent_id": parentID,
			"name":      name,
			"code":      code,
			"category":  category,
			"config":    defaultJSON(row["config"]),
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed org_units %s: %w", code, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableOrgUnits, scopedLookupKey(row["tenant_code"], code), id)
	}

	return nil
}

func (s *Seed) seedUsers(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableUsers]
	cfg := upsertConfig{
		table:              seedTableUsers,
		columns:            []string{"tenant_id", "company_id", "name", "username", "email", "password", "email_verified_at"},
		matchColumns:       []string{"tenant_id", "email"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		email, err := requiredCSVValue(row, "email")
		if err != nil {
			return fmt.Errorf("seed users: %w", err)
		}
		name, err := requiredCSVValue(row, "name")
		if err != nil {
			return fmt.Errorf("seed users %s: %w", email, err)
		}
		username, err := requiredCSVValue(row, "username")
		if err != nil {
			return fmt.Errorf("seed users %s: %w", email, err)
		}
		passwordPlain, err := requiredCSVValue(row, "password_plain")
		if err != nil {
			return fmt.Errorf("seed users %s: %w", email, err)
		}

		tenantID, err := resolveTenantID(state, row["tenant_code"])
		if err != nil {
			return fmt.Errorf("seed users %s: %w", email, err)
		}
		companyScope := scopedLookupScope(row["company_tenant_code"], row["tenant_code"])
		companyID, err := resolveScopedNullableID(state, seedTableOrgUnits, companyScope, row["company_code"], "company")
		if err != nil {
			return fmt.Errorf("seed users %s: %w", email, err)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordPlain), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("seed users %s: hash password: %w", email, err)
		}
		emailVerifiedAt, err := optionalTimeValue(row, "email_verified_at")
		if err != nil {
			return fmt.Errorf("seed users %s: %w", email, err)
		}

		values := map[string]any{
			"tenant_id":         tenantID,
			"company_id":        companyID,
			"name":              name,
			"username":          username,
			"email":             strings.ToLower(strings.TrimSpace(email)),
			"password":          string(hashedPassword),
			"email_verified_at": emailVerifiedAt,
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed users %s: %w", email, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableUsers, scopedLookupKey(row["tenant_code"], email), id)
	}

	return nil
}

func (s *Seed) seedUserRoles(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableUserRoles]

	for _, row := range file.rows {
		email, err := requiredCSVValue(row, "user_email")
		if err != nil {
			return fmt.Errorf("seed user_roles: %w", err)
		}
		roleName, err := requiredCSVValue(row, "role_name")
		if err != nil {
			return fmt.Errorf("seed user_roles %s/%s: %w", email, roleName, err)
		}

		userID, err := resolveUserID(state, row["tenant_code"], email)
		if err != nil {
			return fmt.Errorf("seed user_roles %s/%s: %w", email, roleName, err)
		}
		roleID, err := resolveScopedID(state, seedTableRoles, row["tenant_code"], roleName, "role")
		if err != nil {
			return fmt.Errorf("seed user_roles %s/%s: %w", email, roleName, err)
		}

		if strings.TrimSpace(row["id"]) == "" {
			id, err := newUUIDv7()
			if err != nil {
				return fmt.Errorf("seed user_roles %s/%s: %w", email, roleName, err)
			}
			row["id"] = id
			file.modified = true
		} else if err := validateUUIDv7(row["id"]); err != nil {
			return fmt.Errorf("seed user_roles %s/%s: %w", email, roleName, err)
		}

		query := `
			INSERT INTO user_roles (user_id, role_id)
			VALUES (?, ?)
			ON CONFLICT (user_id, role_id) DO NOTHING
		`
		if _, err := tx.Exec(tx.Rebind(query), *userID, roleID); err != nil {
			return fmt.Errorf("seed user_roles %s/%s: %w", email, roleName, err)
		}
	}

	return nil
}

func (s *Seed) seedInternalUsers(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableInternalUsers]
	cfg := upsertConfig{
		table:              seedTableInternalUsers,
		columns:            []string{"email", "full_name", "password_hash", "role_code", "status"},
		matchColumns:       []string{"email"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		email, err := requiredCSVValue(row, "email")
		if err != nil {
			return fmt.Errorf("seed internal_users: %w", err)
		}
		fullName, err := requiredCSVValue(row, "full_name")
		if err != nil {
			return fmt.Errorf("seed internal_users %s: %w", email, err)
		}
		passwordPlain, err := requiredCSVValue(row, "password_plain")
		if err != nil {
			return fmt.Errorf("seed internal_users %s: %w", email, err)
		}
		roleCode, err := requiredCSVValue(row, "role_code")
		if err != nil {
			return fmt.Errorf("seed internal_users %s: %w", email, err)
		}
		status, err := requiredCSVValue(row, "status")
		if err != nil {
			return fmt.Errorf("seed internal_users %s: %w", email, err)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordPlain), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("seed internal_users %s: hash password: %w", email, err)
		}

		values := map[string]any{
			"email":         strings.ToLower(strings.TrimSpace(email)),
			"full_name":     fullName,
			"password_hash": string(hashedPassword),
			"role_code":     roleCode,
			"status":        status,
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed internal_users %s: %w", email, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableInternalUsers, tenantLookupKey(email), id)
	}

	return nil
}

func (s *Seed) seedInternalCurrencies(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableInternalCurrencies]

	defaultCode := ""
	for _, row := range file.rows {
		code, err := requiredCSVValue(row, "code")
		if err != nil {
			return fmt.Errorf("seed internal_currencies: %w", err)
		}
		isDefault, err := requiredBoolValue(row, "is_default")
		if err != nil {
			return fmt.Errorf("seed internal_currencies %s: %w", code, err)
		}
		if !isDefault {
			continue
		}
		if defaultCode != "" {
			return fmt.Errorf("seed internal_currencies: multiple default currencies defined: %s and %s", defaultCode, code)
		}
		defaultCode = code
	}

	if defaultCode != "" {
		query := `
			UPDATE internal_currencies
			SET
				is_default = FALSE,
				updated_at = CURRENT_TIMESTAMP
			WHERE deleted_at IS NULL
				AND code <> ?
				AND is_default = TRUE
		`
		if _, err := tx.Exec(tx.Rebind(query), defaultCode); err != nil {
			return fmt.Errorf("seed internal_currencies: clear previous default: %w", err)
		}
	}

	for _, row := range file.rows {
		code, err := requiredCSVValue(row, "code")
		if err != nil {
			return fmt.Errorf("seed internal_currencies: %w", err)
		}
		name, err := requiredCSVValue(row, "name")
		if err != nil {
			return fmt.Errorf("seed internal_currencies %s: %w", code, err)
		}
		symbol, err := requiredCSVValue(row, "symbol")
		if err != nil {
			return fmt.Errorf("seed internal_currencies %s: %w", code, err)
		}
		decimalPlaces, err := requiredIntValue(row, "decimal_places")
		if err != nil {
			return fmt.Errorf("seed internal_currencies %s: %w", code, err)
		}
		isActive, err := requiredBoolValue(row, "is_active")
		if err != nil {
			return fmt.Errorf("seed internal_currencies %s: %w", code, err)
		}
		isDefault, err := requiredBoolValue(row, "is_default")
		if err != nil {
			return fmt.Errorf("seed internal_currencies %s: %w", code, err)
		}
		sortOrder, err := requiredIntValue(row, "sort_order")
		if err != nil {
			return fmt.Errorf("seed internal_currencies %s: %w", code, err)
		}

		query := `
			INSERT INTO internal_currencies (
				code,
				name,
				symbol,
				decimal_places,
				is_active,
				is_default,
				sort_order,
				metadata
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (code) DO UPDATE SET
				name = EXCLUDED.name,
				symbol = EXCLUDED.symbol,
				decimal_places = EXCLUDED.decimal_places,
				is_active = EXCLUDED.is_active,
				is_default = EXCLUDED.is_default,
				sort_order = EXCLUDED.sort_order,
				metadata = EXCLUDED.metadata,
				updated_at = CURRENT_TIMESTAMP,
				deleted_at = NULL
		`
		if _, err := tx.Exec(
			tx.Rebind(query),
			code,
			name,
			symbol,
			decimalPlaces,
			isActive,
			isDefault,
			sortOrder,
			defaultJSON(row["metadata"]),
		); err != nil {
			return fmt.Errorf("seed internal_currencies %s: %w", code, err)
		}
	}

	return nil
}

func (s *Seed) seedInternalPaymentProviderCurrencies(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableInternalProviderCurrencies]

	for _, row := range file.rows {
		provider, err := requiredCSVValue(row, "provider")
		if err != nil {
			return fmt.Errorf("seed internal_payment_provider_currencies: %w", err)
		}
		currencyCode, err := requiredCSVValue(row, "currency_code")
		if err != nil {
			return fmt.Errorf("seed internal_payment_provider_currencies %s: %w", provider, err)
		}
		if _, err := resolveID(state, seedTableInternalCurrencies, currencyCode, "internal currency"); err != nil {
			return fmt.Errorf("seed internal_payment_provider_currencies %s/%s: %w", provider, currencyCode, err)
		}
		isActive, err := requiredBoolValue(row, "is_active")
		if err != nil {
			return fmt.Errorf("seed internal_payment_provider_currencies %s/%s: %w", provider, currencyCode, err)
		}

		query := `
			INSERT INTO internal_payment_provider_currencies (
				provider,
				currency_code,
				is_active,
				min_amount,
				max_amount,
				metadata
			) VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT (provider, currency_code) DO UPDATE SET
				is_active = EXCLUDED.is_active,
				min_amount = EXCLUDED.min_amount,
				max_amount = EXCLUDED.max_amount,
				metadata = EXCLUDED.metadata,
				updated_at = CURRENT_TIMESTAMP,
				deleted_at = NULL
		`
		if _, err := tx.Exec(
			tx.Rebind(query),
			provider,
			currencyCode,
			isActive,
			nullableStringValue(row["min_amount"]),
			nullableStringValue(row["max_amount"]),
			defaultJSON(row["metadata"]),
		); err != nil {
			return fmt.Errorf("seed internal_payment_provider_currencies %s/%s: %w", provider, currencyCode, err)
		}
	}

	return nil
}

func (s *Seed) seedInternalProducts(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableInternalProducts]
	cfg := upsertConfig{
		table:              seedTableInternalProducts,
		columns:            []string{"code", "name", "description", "status", "metadata"},
		matchColumns:       []string{"code"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		code, err := requiredCSVValue(row, "code")
		if err != nil {
			return fmt.Errorf("seed internal_products: %w", err)
		}
		name, err := requiredCSVValue(row, "name")
		if err != nil {
			return fmt.Errorf("seed internal_products %s: %w", code, err)
		}
		status, err := requiredCSVValue(row, "status")
		if err != nil {
			return fmt.Errorf("seed internal_products %s: %w", code, err)
		}
		values := map[string]any{
			"code":        code,
			"name":        name,
			"description": strings.TrimSpace(row["description"]),
			"status":      status,
			"metadata":    defaultJSON(row["metadata"]),
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed internal_products %s: %w", code, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableInternalProducts, code, id)
	}

	return nil
}

func (s *Seed) seedInternalProductPricings(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableInternalProductPricings]
	cfg := upsertConfig{
		table:              seedTableInternalProductPricings,
		columns:            []string{"internal_product_id", "code", "name", "description", "status", "metadata"},
		matchColumns:       []string{"internal_product_id", "code"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		code, err := requiredCSVValue(row, "code")
		if err != nil {
			return fmt.Errorf("seed internal_product_pricings: %w", err)
		}
		name, err := requiredCSVValue(row, "name")
		if err != nil {
			return fmt.Errorf("seed internal_product_pricings %s: %w", code, err)
		}
		status, err := requiredCSVValue(row, "status")
		if err != nil {
			return fmt.Errorf("seed internal_product_pricings %s: %w", code, err)
		}
		productCode, err := requiredCSVValue(row, "product_code")
		if err != nil {
			return fmt.Errorf("seed internal_product_pricings %s: %w", code, err)
		}
		productID, err := resolveID(state, seedTableInternalProducts, productCode, "internal product")
		if err != nil {
			return fmt.Errorf("seed internal_product_pricings %s: %w", code, err)
		}

		values := map[string]any{
			"internal_product_id": productID,
			"code":                code,
			"name":                name,
			"description":         strings.TrimSpace(row["description"]),
			"status":              status,
			"metadata":            defaultJSON(row["metadata"]),
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed internal_product_pricings %s: %w", code, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableInternalProductPricings, productCode+"::"+code, id)
	}

	return nil
}

func (s *Seed) seedInternalProductPrices(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableInternalProductPrices]
	cfg := upsertConfig{
		table:              seedTableInternalProductPrices,
		columns:            []string{"internal_product_pricing_id", "currency_code", "amount", "started_at", "ended_at", "metadata"},
		matchColumns:       []string{"internal_product_pricing_id", "currency_code", "started_at"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		productCode, err := requiredCSVValue(row, "product_code")
		if err != nil {
			return fmt.Errorf("seed internal_product_prices: %w", err)
		}
		pricingCode, err := requiredCSVValue(row, "pricing_code")
		if err != nil {
			return fmt.Errorf("seed internal_product_prices %s/%s: %w", productCode, pricingCode, err)
		}
		currencyCode, err := requiredCSVValue(row, "currency_code")
		if err != nil {
			return fmt.Errorf("seed internal_product_prices %s/%s: %w", productCode, pricingCode, err)
		}
		amount, err := requiredCSVValue(row, "amount")
		if err != nil {
			return fmt.Errorf("seed internal_product_prices %s/%s: %w", productCode, pricingCode, err)
		}
		startedAt, err := requiredCSVValue(row, "started_at")
		if err != nil {
			return fmt.Errorf("seed internal_product_prices %s/%s: %w", productCode, pricingCode, err)
		}
		pricingID, err := resolveID(state, seedTableInternalProductPricings, productCode+"::"+pricingCode, "internal product pricing")
		if err != nil {
			return fmt.Errorf("seed internal_product_prices %s/%s: %w", productCode, pricingCode, err)
		}
		endedAt, err := optionalTimeValue(row, "ended_at")
		if err != nil {
			return fmt.Errorf("seed internal_product_prices %s/%s: %w", productCode, pricingCode, err)
		}
		startedAtTime, err := optionalTimeValue(row, "started_at")
		if err != nil || startedAtTime == nil {
			if err == nil {
				err = fmt.Errorf("started_at is required")
			}
			return fmt.Errorf("seed internal_product_prices %s/%s: %w", productCode, pricingCode, err)
		}

		values := map[string]any{
			"internal_product_pricing_id": pricingID,
			"currency_code":               currencyCode,
			"amount":                      amount,
			"started_at":                  startedAtTime,
			"ended_at":                    endedAt,
			"metadata":                    defaultJSON(row["metadata"]),
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed internal_product_prices %s/%s: %w", productCode, pricingCode, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableInternalProductPrices, productCode+"::"+pricingCode+"::"+currencyCode+"::"+startedAt, id)
	}

	return nil
}

func (s *Seed) seedJobPositions(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableJobPositions]
	cfg := upsertConfig{
		table:              seedTableJobPositions,
		columns:            []string{"tenant_id", "name", "grade"},
		matchColumns:       []string{"tenant_id", "name"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		name, err := requiredCSVValue(row, "name")
		if err != nil {
			return fmt.Errorf("seed job_positions: %w", err)
		}
		tenantID, err := resolveTenantID(state, row["tenant_code"])
		if err != nil {
			return fmt.Errorf("seed job_positions %s: %w", name, err)
		}

		values := map[string]any{
			"tenant_id": tenantID,
			"name":      name,
			"grade":     nullableStringValue(row["grade"]),
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed job_positions %s: %w", name, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableJobPositions, scopedLookupKey(row["tenant_code"], name), id)
	}

	return nil
}

func (s *Seed) seedWorkLocations(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableWorkLocations]
	cfg := upsertConfig{
		table:              seedTableWorkLocations,
		columns:            []string{"tenant_id", "org_unit_id", "name", "address", "latitude", "longitude", "radius_meters"},
		matchColumns:       []string{"tenant_id", "name"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		name, err := requiredCSVValue(row, "name")
		if err != nil {
			return fmt.Errorf("seed work_locations: %w", err)
		}
		tenantID, err := resolveTenantID(state, row["tenant_code"])
		if err != nil {
			return fmt.Errorf("seed work_locations %s: %w", name, err)
		}
		orgUnitScope := scopedLookupScope(row["org_unit_tenant_code"], row["tenant_code"])
		orgUnitID, err := resolveScopedNullableID(state, seedTableOrgUnits, orgUnitScope, row["org_unit_code"], "org unit")
		if err != nil {
			return fmt.Errorf("seed work_locations %s: %w", name, err)
		}
		latitude, err := optionalFloatValue(row, "latitude")
		if err != nil {
			return fmt.Errorf("seed work_locations %s: %w", name, err)
		}
		longitude, err := optionalFloatValue(row, "longitude")
		if err != nil {
			return fmt.Errorf("seed work_locations %s: %w", name, err)
		}
		radiusMeters, err := optionalIntValue(row, "radius_meters")
		if err != nil {
			return fmt.Errorf("seed work_locations %s: %w", name, err)
		}

		values := map[string]any{
			"tenant_id":     tenantID,
			"org_unit_id":   orgUnitID,
			"name":          name,
			"address":       nullableStringValue(row["address"]),
			"latitude":      latitude,
			"longitude":     longitude,
			"radius_meters": radiusMeters,
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed work_locations %s: %w", name, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableWorkLocations, scopedLookupKey(row["tenant_code"], name), id)
	}

	return nil
}

func (s *Seed) seedWorkShifts(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableWorkShifts]
	cfg := upsertConfig{
		table:              seedTableWorkShifts,
		columns:            []string{"tenant_id", "name", "timezone", "start_time", "end_time", "grace_period_minutes"},
		matchColumns:       []string{"tenant_id", "name"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		name, err := requiredCSVValue(row, "name")
		if err != nil {
			return fmt.Errorf("seed work_shifts: %w", err)
		}
		timezone, err := requiredCSVValue(row, "timezone")
		if err != nil {
			return fmt.Errorf("seed work_shifts %s: %w", name, err)
		}
		tenantID, err := resolveTenantID(state, row["tenant_code"])
		if err != nil {
			return fmt.Errorf("seed work_shifts %s: %w", name, err)
		}
		gracePeriodMinutes, err := requiredIntValue(row, "grace_period_minutes")
		if err != nil {
			return fmt.Errorf("seed work_shifts %s: %w", name, err)
		}

		values := map[string]any{
			"tenant_id":            tenantID,
			"name":                 name,
			"timezone":             timezone,
			"start_time":           nullableStringValue(row["start_time"]),
			"end_time":             nullableStringValue(row["end_time"]),
			"grace_period_minutes": gracePeriodMinutes,
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed work_shifts %s: %w", name, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableWorkShifts, scopedLookupKey(row["tenant_code"], name), id)
	}

	return nil
}

func (s *Seed) seedEmployees(tx *sqlx.Tx, state *csvSeedState) error {
	file := state.files[seedTableEmployees]
	cfg := upsertConfig{
		table:              seedTableEmployees,
		columns:            []string{"tenant_id", "employee_no", "full_name", "user_id", "org_unit_id", "job_position_id", "location_id", "shift_id", "email", "status", "join_date", "join_date_timezone"},
		matchColumns:       []string{"tenant_id", "employee_no"},
		hasUpdatedAt:       true,
		supportsSoftDelete: true,
	}

	for _, row := range file.rows {
		employeeNo, err := requiredCSVValue(row, "employee_no")
		if err != nil {
			return fmt.Errorf("seed employees: %w", err)
		}
		fullName, err := requiredCSVValue(row, "full_name")
		if err != nil {
			return fmt.Errorf("seed employees %s: %w", employeeNo, err)
		}
		status, err := requiredCSVValue(row, "status")
		if err != nil {
			return fmt.Errorf("seed employees %s: %w", employeeNo, err)
		}

		tenantID, err := resolveTenantID(state, row["tenant_code"])
		if err != nil {
			return fmt.Errorf("seed employees %s: %w", employeeNo, err)
		}
		userID, err := resolveUserID(state, row["tenant_code"], row["user_email"])
		if err != nil {
			return fmt.Errorf("seed employees %s: %w", employeeNo, err)
		}
		orgUnitScope := scopedLookupScope(row["org_unit_tenant_code"], row["tenant_code"])
		orgUnitID, err := resolveScopedNullableID(state, seedTableOrgUnits, orgUnitScope, row["org_unit_code"], "org unit")
		if err != nil {
			return fmt.Errorf("seed employees %s: %w", employeeNo, err)
		}
		jobPositionScope := scopedLookupScope(row["job_position_tenant_code"], row["tenant_code"])
		jobPositionID, err := resolveScopedNullableID(state, seedTableJobPositions, jobPositionScope, row["job_position_name"], "job position")
		if err != nil {
			return fmt.Errorf("seed employees %s: %w", employeeNo, err)
		}
		locationScope := scopedLookupScope(row["location_tenant_code"], row["tenant_code"])
		locationID, err := resolveScopedNullableID(state, seedTableWorkLocations, locationScope, row["location_name"], "work location")
		if err != nil {
			return fmt.Errorf("seed employees %s: %w", employeeNo, err)
		}
		shiftScope := scopedLookupScope(row["shift_tenant_code"], row["tenant_code"])
		shiftID, err := resolveScopedID(state, seedTableWorkShifts, shiftScope, row["shift_name"], "work shift")
		if err != nil {
			return fmt.Errorf("seed employees %s: %w", employeeNo, err)
		}

		values := map[string]any{
			"tenant_id":          tenantID,
			"employee_no":        employeeNo,
			"full_name":          fullName,
			"user_id":            userID,
			"org_unit_id":        orgUnitID,
			"job_position_id":    jobPositionID,
			"location_id":        locationID,
			"shift_id":           &shiftID,
			"email":              nullableStringValue(row["email"]),
			"status":             status,
			"join_date":          nullableStringValue(row["join_date"]),
			"join_date_timezone": nullableStringValue(row["join_date_timezone"]),
		}

		id, backfilled, err := upsertRecord(tx, cfg, row["id"], values)
		if err != nil {
			return fmt.Errorf("seed employees %s: %w", employeeNo, err)
		}
		if backfilled {
			row["id"] = id
			file.modified = true
		}
		state.registerID(seedTableEmployees, scopedLookupKey(row["tenant_code"], employeeNo), id)
	}

	return nil
}

func resolveTenantID(state *csvSeedState, tenantCode string) (*string, error) {
	tenantCode = strings.TrimSpace(tenantCode)
	if tenantCode == "" {
		return nil, nil
	}
	id, ok := state.lookupID(seedTableTenants, tenantLookupKey(tenantCode))
	if !ok {
		return nil, fmt.Errorf("tenant with code %q not found", tenantCode)
	}
	return &id, nil
}

func resolveUserID(state *csvSeedState, tenantCode string, email string) (*string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil, nil
	}

	if strings.TrimSpace(tenantCode) != "" {
		id, ok := state.lookupID(seedTableUsers, scopedLookupKey(tenantCode, email))
		if !ok {
			return nil, fmt.Errorf("user %q not found in tenant %q", email, tenantCode)
		}
		return &id, nil
	}

	entries, ok := state.ids[seedTableUsers]
	if !ok {
		return nil, fmt.Errorf("user %q not found", email)
	}

	var matchedID string
	matches := 0
	suffix := "::" + email
	for key, id := range entries {
		if strings.HasSuffix(key, suffix) {
			matchedID = id
			matches++
		}
	}

	if matches == 0 {
		return nil, fmt.Errorf("user %q not found", email)
	}
	if matches > 1 {
		return nil, fmt.Errorf("user %q is ambiguous across tenants", email)
	}

	return &matchedID, nil
}

func resolveScopedID(state *csvSeedState, tableName string, scope string, value string, label string) (string, error) {
	id, ok := state.lookupID(tableName, scopedLookupKey(scope, value))
	if !ok {
		return "", fmt.Errorf("%s %q not found", label, value)
	}
	return id, nil
}

func resolveScopedNullableID(state *csvSeedState, tableName string, scope string, value string, label string) (*string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	id, err := resolveScopedID(state, tableName, scope, value, label)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func resolveID(state *csvSeedState, tableName string, value string, label string) (string, error) {
	value = strings.TrimSpace(value)
	id, ok := state.lookupID(tableName, value)
	if !ok {
		return "", fmt.Errorf("%s %q not found", label, value)
	}
	return id, nil
}

func resolveTemplateID(state *csvSeedState, tableName string, value string, label string) (string, error) {
	return resolveID(state, tableName, value, label)
}

func scopedLookupScope(primary string, fallback string) string {
	primary = strings.TrimSpace(primary)
	if primary != "" {
		return primary
	}
	return strings.TrimSpace(fallback)
}

func defaultJSON(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "{}"
	}
	return value
}

func nullableStringValue(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
