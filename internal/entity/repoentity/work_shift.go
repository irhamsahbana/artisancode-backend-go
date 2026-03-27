package repoentity

type WorkShift struct {
	ID                 string `db:"id"`
	TenantID           string `db:"tenant_id"`
	Name               string `db:"name"`
	Timezone           string `db:"timezone"`
	StartTime          string `db:"start_time"`
	EndTime            string `db:"end_time"`
	GracePeriodMinutes int    `db:"grace_period_minutes"`
}
