package repoentity

type JobPosition struct {
	ID       string  `db:"id"`
	TenantID string  `db:"tenant_id"`
	Name     string  `db:"name"`
	Grade    *string `db:"grade"`
}
