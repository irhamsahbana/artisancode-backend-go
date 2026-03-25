package repoentity

type UserLogin struct {
	ID          string   `db:"id"`
	Email       string   `db:"email"`
	Password    string   `db:"password"`
	RoleNames   []string `db:"role_names"`
	TenantID    string   `db:"tenant_id"`
	TenantName  string   `db:"tenant_name"`
	UserName    string   `db:"username"`
	CompanyID   *string  `db:"company_id"`
	CompanyName *string  `db:"company_name"`
}
