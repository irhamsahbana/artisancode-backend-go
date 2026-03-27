package repoentity

type OrgUnit struct {
	ID       string  `db:"id"`
	TenantID string  `db:"tenant_id"`
	Code     string  `db:"code"`
	ParentID *string `db:"parent_id"`
	Name     string  `db:"name"`
	Category string  `db:"category"`
}
