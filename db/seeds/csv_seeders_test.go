package seeds

import (
	"strings"
	"testing"
)

func TestResolveScopedIDRequiresTenantScopedRecord(t *testing.T) {
	state := &csvSeedState{
		ids: map[string]map[string]string{
			seedTableRoles: {
				scopedLookupKey("", "owner"): "global-owner-id",
			},
		},
	}

	_, err := resolveScopedID(state, seedTableRoles, "BERUA", "owner", "role")
	if err == nil {
		t.Fatal("resolveScopedID expected missing tenant-scoped role error, got nil")
	}
}

func TestScopedLookupScopeUsesFallback(t *testing.T) {
	scope := scopedLookupScope("", "BERUA")
	if scope != "BERUA" {
		t.Fatalf("scopedLookupScope returned %q, want %q", scope, "BERUA")
	}
}

func TestAllWithTemplateSeedsTenantFirst(t *testing.T) {
	order := seedGroups["all_with_template"]
	if len(order) == 0 {
		t.Fatal("all_with_template order is empty")
	}
	if order[0] != seedTableTenants {
		t.Fatalf("all_with_template first table = %q, want %q", order[0], seedTableTenants)
	}
}

func TestEmployeeSeedGroupIncludesEmployees(t *testing.T) {
	order := seedGroups["employee"]
	if len(order) == 0 {
		t.Fatal("employee order is empty")
	}
	if order[len(order)-1] != seedTableEmployees {
		t.Fatalf("employee last table = %q, want %q", order[len(order)-1], seedTableEmployees)
	}
}

func TestInternalCatalogSeedsCurrenciesBeforePrices(t *testing.T) {
	order := seedGroups["internal_catalog"]
	if len(order) < 5 {
		t.Fatalf("internal_catalog order is too short: %#v", order)
	}
	if order[0] != seedTableInternalCurrencies {
		t.Fatalf("internal_catalog first table = %q, want %q", order[0], seedTableInternalCurrencies)
	}
	if order[1] != seedTableInternalProviderCurrencies {
		t.Fatalf("internal_catalog second table = %q, want %q", order[1], seedTableInternalProviderCurrencies)
	}
	if order[len(order)-1] != seedTableInternalProductPrices {
		t.Fatalf("internal_catalog last table = %q, want %q", order[len(order)-1], seedTableInternalProductPrices)
	}
}

func TestResolveUserIDByEmail(t *testing.T) {
	state := &csvSeedState{
		ids: map[string]map[string]string{
			seedTableUsers: {
				scopedLookupKey("BERUA", "attendance.employee.001@example.com"): "user-001",
			},
		},
	}

	id, err := resolveUserID(state, "BERUA", "attendance.employee.001@example.com")
	if err != nil {
		t.Fatalf("resolveUserID returned error: %v", err)
	}
	if id == nil || *id != "user-001" {
		t.Fatalf("resolveUserID returned %v, want %q", id, "user-001")
	}
}

func TestResolveUserIDFallsBackToUniqueEmailWhenTenantEmpty(t *testing.T) {
	state := &csvSeedState{
		ids: map[string]map[string]string{
			seedTableUsers: {
				scopedLookupKey("BERUA", "beruang@beruang.com"): "owner-001",
			},
		},
	}

	id, err := resolveUserID(state, "", "beruang@beruang.com")
	if err != nil {
		t.Fatalf("resolveUserID returned error: %v", err)
	}
	if id == nil || *id != "owner-001" {
		t.Fatalf("resolveUserID returned %v, want %q", id, "owner-001")
	}
}

func TestResolveUserIDReturnsAmbiguousWhenEmailExistsInMultipleTenants(t *testing.T) {
	state := &csvSeedState{
		ids: map[string]map[string]string{
			seedTableUsers: {
				scopedLookupKey("BERUA", "shared@example.com"): "user-001",
				scopedLookupKey("KOPI", "shared@example.com"):  "user-002",
			},
		},
	}

	_, err := resolveUserID(state, "", "shared@example.com")
	if err == nil {
		t.Fatal("resolveUserID expected ambiguous error, got nil")
	}
	if !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("resolveUserID error = %q, want ambiguous", err.Error())
	}
}

func TestTenantScopedRolePermissionsHaveScopedRolesAndPermissions(t *testing.T) {
	state, err := newCSVSeedState()
	if err != nil {
		t.Fatalf("newCSVSeedState returned error: %v", err)
	}

	for _, row := range state.files[seedTableRolePerms].rows {
		scope := strings.TrimSpace(row["tenant_code"])
		if scope == "" {
			continue
		}

		roleKey := scopedLookupKey(scope, row["role_name"])
		if _, ok := state.lookupID(seedTableRoles, roleKey); !ok {
			t.Fatalf("tenant-scoped role %q for %q is missing", row["role_name"], scope)
		}

		permissionKey := scopedLookupKey(scope, row["permission_name"])
		if _, ok := state.lookupID(seedTablePermissions, permissionKey); !ok {
			t.Fatalf("tenant-scoped permission %q for %q is missing", row["permission_name"], scope)
		}
	}
}

func TestRegularRBACSeedsAreTenantScoped(t *testing.T) {
	state, err := newCSVSeedState()
	if err != nil {
		t.Fatalf("newCSVSeedState returned error: %v", err)
	}

	for _, tableName := range []string{seedTableRoles, seedTablePermissions, seedTableRolePerms} {
		for _, row := range state.files[tableName].rows {
			if strings.TrimSpace(row["tenant_code"]) == "" {
				t.Fatalf("%s contains non-tenant-scoped row: %#v", tableName, row)
			}
		}
	}
}

func TestInternalTemplateRolePermissionsHaveTemplateRolesAndPermissions(t *testing.T) {
	state, err := newCSVSeedState()
	if err != nil {
		t.Fatalf("newCSVSeedState returned error: %v", err)
	}

	for _, row := range state.files[seedTableInternalTemplateRolePerms].rows {
		if _, ok := state.lookupID(seedTableInternalTemplateRoles, row["role_name"]); !ok {
			t.Fatalf("template role permission role %q is missing", row["role_name"])
		}
		if _, ok := state.lookupID(seedTableInternalTemplatePerms, row["permission_name"]); !ok {
			t.Fatalf("template role permission permission %q is missing", row["permission_name"])
		}
	}
}

func TestTenantScopedUserRolesHaveScopedRoles(t *testing.T) {
	state, err := newCSVSeedState()
	if err != nil {
		t.Fatalf("newCSVSeedState returned error: %v", err)
	}

	for _, row := range state.files[seedTableUserRoles].rows {
		scope := strings.TrimSpace(row["tenant_code"])
		if scope == "" {
			t.Fatalf("user role %q/%q is missing tenant_code", row["user_email"], row["role_name"])
		}

		roleKey := scopedLookupKey(scope, row["role_name"])
		if _, ok := state.lookupID(seedTableRoles, roleKey); !ok {
			t.Fatalf("tenant-scoped user role %q for %q is missing", row["role_name"], scope)
		}
	}
}

func TestEmployeeSeedsRequireShiftName(t *testing.T) {
	state, err := newCSVSeedState()
	if err != nil {
		t.Fatalf("newCSVSeedState returned error: %v", err)
	}

	for _, row := range state.files[seedTableEmployees].rows {
		if strings.TrimSpace(row["shift_name"]) == "" {
			t.Fatalf("employee %q is missing shift_name", row["employee_no"])
		}
	}
}

func TestEmployeeUsersHaveEmailVerifiedAt(t *testing.T) {
	state, err := newCSVSeedState()
	if err != nil {
		t.Fatalf("newCSVSeedState returned error: %v", err)
	}

	for _, row := range state.files[seedTableUsers].rows {
		email := strings.TrimSpace(row["email"])
		if !strings.HasPrefix(email, "emp") {
			continue
		}
		if strings.TrimSpace(row["email_verified_at"]) == "" {
			t.Fatalf("employee user %q is missing email_verified_at", email)
		}
	}
}

func TestOwnerSeedHasEmailVerifiedAt(t *testing.T) {
	state, err := newCSVSeedState()
	if err != nil {
		t.Fatalf("newCSVSeedState returned error: %v", err)
	}

	for _, row := range state.files[seedTableUsers].rows {
		email := strings.TrimSpace(row["email"])
		if email != "beruang@beruang.com" {
			continue
		}

		if strings.TrimSpace(row["email_verified_at"]) == "" {
			t.Fatalf("owner user %q is missing email_verified_at", email)
		}

		return
	}

	t.Fatal("owner seed user beruang@beruang.com is missing")
}

func TestInternalProviderCurrencySeedsReferenceKnownCurrencies(t *testing.T) {
	state, err := newCSVSeedState()
	if err != nil {
		t.Fatalf("newCSVSeedState returned error: %v", err)
	}

	for _, row := range state.files[seedTableInternalProviderCurrencies].rows {
		currencyCode := strings.TrimSpace(row["currency_code"])
		if currencyCode == "" {
			t.Fatalf("provider currency row is missing currency_code: %#v", row)
		}
		if _, ok := state.lookupID(seedTableInternalCurrencies, currencyCode); !ok {
			t.Fatalf("provider currency %q references unknown currency %q", row["provider"], currencyCode)
		}
	}
}
