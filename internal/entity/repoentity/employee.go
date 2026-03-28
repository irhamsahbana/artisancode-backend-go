package repoentity

type Employee struct {
	ID            string  `db:"id"`
	TenantID      string  `db:"tenant_id"`
	EmployeeNo    string  `db:"employee_no"`
	FullName      string  `db:"full_name"`
	Email         *string `db:"email"`
	OrgUnitID     *string `db:"org_unit_id"`
	JobPositionID *string `db:"job_position_id"`
	LocationID    *string `db:"location_id"`
	ShiftID       *string `db:"shift_id"`
	Status        string  `db:"status"`
	JoinDate      *string `db:"join_date"`
}
