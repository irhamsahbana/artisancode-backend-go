package repoentity

type WorkLocation struct {
	ID           string   `db:"id"`
	TenantID     string   `db:"tenant_id"`
	OrgUnitID    *string  `db:"org_unit_id"`
	OrgUnitName  *string  `db:"org_unit_name"`
	Name         string   `db:"name"`
	Address      *string  `db:"address"`
	Timezone     string   `db:"timezone"`
	Latitude     *float64 `db:"latitude"`
	Longitude    *float64 `db:"longitude"`
	RadiusMeters *int     `db:"radius_meters"`
}