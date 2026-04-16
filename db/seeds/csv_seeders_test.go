package seeds

import "testing"

func TestResolveScopedIDFallsBackToGlobal(t *testing.T) {
	state := &csvSeedState{
		ids: map[string]map[string]string{
			seedTableRoles: {
				scopedLookupKey("", "owner"): "global-owner-id",
			},
		},
	}

	id, err := resolveScopedID(state, seedTableRoles, "BERUA", "owner", "role")
	if err != nil {
		t.Fatalf("resolveScopedID returned error: %v", err)
	}
	if id != "global-owner-id" {
		t.Fatalf("resolveScopedID returned %q, want %q", id, "global-owner-id")
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

func TestResolveUserIDByEmail(t *testing.T) {
	state := &csvSeedState{
		ids: map[string]map[string]string{
			seedTableUsers: {
				"attendance.employee.001@example.com": "user-001",
			},
		},
	}

	id, err := resolveUserID(state, "attendance.employee.001@example.com")
	if err != nil {
		t.Fatalf("resolveUserID returned error: %v", err)
	}
	if id == nil || *id != "user-001" {
		t.Fatalf("resolveUserID returned %v, want %q", id, "user-001")
	}
}
